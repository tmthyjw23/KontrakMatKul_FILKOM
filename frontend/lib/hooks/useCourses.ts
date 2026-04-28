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
  return {
    id: String(course.id),
    code: course.code,
    name: course.name,
    sks: Number(course.sks),
    lecturer: course.lecturer ?? "Dosen belum ditentukan",
    schedules: (course.schedules ?? []).map(normalizeSchedule),
    color: course.color ?? fallbackColors[index % fallbackColors.length],
  };
}

async function fetchCourses(): Promise<Course[]> {
  const response = await apiClient.get<
    ApiResponse<BackendCourse[] | { courses: BackendCourse[] }>
  >("/courses");

  const payload = response.data.data;
  const rawCourses = Array.isArray(payload) ? payload : payload.courses ?? [];

  return rawCourses.map(normalizeCourse);
}

export function useCourses() {
  return useQuery({
    queryKey: ["courses"],
    queryFn: fetchCourses,
  });
}
