"use client";

import { create } from "zustand";

import type { Course, Schedule, VisualConflictMap } from "@/src/types/course";

const mockCourses: Course[] = [
  {
    id: "IFAP272",
    code: "IFAP272",
    name: "Automata",
    sks: 3,
    lecturer: "Dr. Maria Kairupan",
    color: "from-emerald-300/18 via-white/8 to-white/5",
    schedules: [
      {
        day: "Senin",
        start_time: "10:10",
        end_time: "11:30",
        room: "GK-302",
      },
      {
        day: "Rabu",
        start_time: "10:10",
        end_time: "11:30",
        room: "GK-302",
      },
    ],
  },
  {
    id: "IFMI252",
    code: "IFMI252",
    name: "Database",
    sks: 3,
    lecturer: "Ir. Jonathan Paat, M.Kom",
    color: "from-white/14 via-white/7 to-emerald-200/10",
    schedules: [
      {
        day: "Selasa",
        start_time: "07:10",
        end_time: "08:30",
        room: "LAB-201",
      },
      {
        day: "Kamis",
        start_time: "07:10",
        end_time: "08:30",
        room: "LAB-201",
      },
    ],
  },
];

const MAX_SKS = 24;

export interface ContractState {
  courses: Course[];
  selectedCourses: Course[];
  totalSks: number;
  maxSks: number;
  addCourse: (courseId: Course["id"]) => void;
  removeCourse: (courseId: Course["id"]) => void;
  getVisualConflictMap: () => VisualConflictMap;
  hasVisualConflict: (courseId: Course["id"]) => boolean;
}

function parseTimeToMinutes(time: string): number | null {
  const [hours, minutes] = time.split(":").map(Number);

  if (
    Number.isNaN(hours) ||
    Number.isNaN(minutes) ||
    hours < 0 ||
    hours > 23 ||
    minutes < 0 ||
    minutes > 59
  ) {
    return null;
  }

  return hours * 60 + minutes;
}

function schedulesOverlap(a: Schedule, b: Schedule): boolean {
  if (a.day !== b.day) {
    return false;
  }

  const startA = parseTimeToMinutes(a.start_time);
  const endA = parseTimeToMinutes(a.end_time);
  const startB = parseTimeToMinutes(b.start_time);
  const endB = parseTimeToMinutes(b.end_time);

  if (
    startA === null ||
    endA === null ||
    startB === null ||
    endB === null
  ) {
    return false;
  }

  return startA < endB && endA > startB;
}

function buildVisualConflictMap(selectedCourses: Course[]): VisualConflictMap {
  const conflictMap: VisualConflictMap = {};

  for (let i = 0; i < selectedCourses.length; i += 1) {
    const currentCourse = selectedCourses[i];

    for (let j = i + 1; j < selectedCourses.length; j += 1) {
      const comparedCourse = selectedCourses[j];

      const hasConflict = currentCourse.schedules.some((currentSchedule) =>
        comparedCourse.schedules.some((comparedSchedule) =>
          schedulesOverlap(currentSchedule, comparedSchedule)
        )
      );

      if (!hasConflict) {
        continue;
      }

      conflictMap[currentCourse.id] = [
        ...(conflictMap[currentCourse.id] ?? []),
        comparedCourse.id,
      ];
      conflictMap[comparedCourse.id] = [
        ...(conflictMap[comparedCourse.id] ?? []),
        currentCourse.id,
      ];
    }
  }

  return conflictMap;
}

export const useContractStore = create<ContractState>((set, get) => ({
  courses: mockCourses,
  selectedCourses: [],
  totalSks: 0,
  maxSks: MAX_SKS,
  addCourse: (courseId) => {
    const { courses, selectedCourses } = get();
    const course = courses.find((item) => item.id === courseId);

    if (!course) {
      return;
    }

    const alreadySelected = selectedCourses.some((item) => item.id === courseId);

    if (alreadySelected) {
      return;
    }

    const nextSelectedCourses = [...selectedCourses, course];

    set({
      selectedCourses: nextSelectedCourses,
      totalSks: nextSelectedCourses.reduce((total, item) => total + item.sks, 0),
    });
  },
  removeCourse: (courseId) => {
    const nextSelectedCourses = get().selectedCourses.filter(
      (item) => item.id !== courseId
    );

    set({
      selectedCourses: nextSelectedCourses,
      totalSks: nextSelectedCourses.reduce((total, item) => total + item.sks, 0),
    });
  },
  getVisualConflictMap: () => buildVisualConflictMap(get().selectedCourses),
  hasVisualConflict: (courseId) => Boolean(buildVisualConflictMap(get().selectedCourses)[courseId]?.length),
}));
