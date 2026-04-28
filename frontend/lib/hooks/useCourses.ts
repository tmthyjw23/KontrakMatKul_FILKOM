"use client";

import { useQuery } from "@tanstack/react-query";

import { apiClient } from "@/lib/api/client";
import type { Course, Schedule } from "@/src/types/course";

type ApiResponse<T> = {
  code: number;
  status: string;
  message: string;
  data: T;
};

type BackendSchedule = Partial<Schedule> & {
  day_of_week?: string;
};

type BackendCourse = {
  id: string | number;
  code: string;
  name: string;
  sks: number | string;
  lecturer?: string | null;
  schedules?: BackendSchedule[];
  color?: string | null;
};

const fallbackColors = [
  "from-emerald-300/18 via-white/8 to-white/5",
  "from-white/14 via-white/7 to-emerald-200/10",
  "from-white/10 via-zinc-200/10 to-emerald-300/12",
  "from-zinc-200/10 via-white/6 to-emerald-100/10",
];

function normalizeTime(value?: string): string {
  if (!value) {
    return "00:00";
  }

  return value.slice(0, 5);
}

function normalizeSchedule(schedule: BackendSchedule): Schedule {
  return {
    day: schedule.day ?? schedule.day_of_week ?? "TBA",
    start_time: normalizeTime(schedule.start_time),
    end_time: normalizeTime(schedule.end_time),
    room: schedule.room ?? "Room TBA",
  };
}

function normalizeCourse(course: BackendCourse, index: number): Course {
  // Extract correct fields based on actual backend response format
  const code = course.code;
  const credits = (course as any).credits ?? course.sks;
  const lecturerName = (course as any).lecturer_name ?? course.lecturer;

  // Generate a mock schedule if the backend doesn't provide one
  const mockDays = ["Senin", "Selasa", "Rabu", "Kamis", "Jumat"];
  const mockStarts = ["08:00", "10:00", "13:00", "15:00"];
  const hasSchedules = course.schedules && course.schedules.length > 0;
  const schedules = hasSchedules 
    ? course.schedules!.map(normalizeSchedule)
    : [{
        day: mockDays[index % mockDays.length],
        start_time: mockStarts[index % mockStarts.length],
        end_time: `${parseInt(mockStarts[index % mockStarts.length].split(":")[0]) + 2}:00`,
        room: `Ruang 10${(index % 9) + 1}`
      }];

  return {
    id: String(course.id ?? code),
    code: code,
    name: course.name,
    sks: Number(credits || 0),
    lecturer: lecturerName ?? "Dosen belum ditentukan",
    schedules: schedules,
    color: course.color ?? fallbackColors[index % fallbackColors.length],
  };
}

async function fetchCourses(): Promise<Course[]> {
  const response = await apiClient.get<unknown>("/courses");

  // Handle both wrapped { code, status, message, data: [...] } (after Antigravity)
  // and unwrapped [...] (current backend state)
  const responseData = response.data as any;
  const payload =
    responseData?.data != null ? responseData.data : responseData;
  const rawCourses: BackendCourse[] = Array.isArray(payload)
    ? payload
    : payload?.courses ?? [];

  return rawCourses.map(normalizeCourse);
}

export function useCourses() {
  return useQuery({
    queryKey: ["courses"],
    queryFn: fetchCourses,
  });
}
