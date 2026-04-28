"use client";

import { create } from "zustand";

import type { Course, Schedule, VisualConflictMap } from "@/src/types/course";

const MAX_SKS = 24;

export interface ContractState {
  courses: Course[];
  selectedCourses: Course[];
  totalSks: number;
  maxSks: number;
  setCourses: (courses: Course[]) => void;
  clearSelectedCourses: () => void;
  addCourse: (courseId: Course["id"]) => void;
  removeCourse: (courseId: Course["id"]) => void;
  getVisualConflictMap: () => VisualConflictMap;
  hasVisualConflict: (courseId: Course["id"]) => boolean;
}

function calculateTotalSks(selectedCourses: Course[]): number {
  return selectedCourses.reduce((total, item) => total + item.sks, 0);
}

function haveSameCourseIds(left: Course[], right: Course[]): boolean {
  if (left.length !== right.length) {
    return false;
  }

  return left.every((course, index) => course.id === right[index]?.id);
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
  courses: [],
  selectedCourses: [],
  totalSks: 0,
  maxSks: MAX_SKS,
  setCourses: (courses) => {
    set((state) => {
      const currentSelectedIds = new Set(
        state.selectedCourses.map((course) => course.id)
      );
      const syncedSelectedCourses = courses.filter((course) =>
        currentSelectedIds.has(course.id)
      );

      if (
        haveSameCourseIds(state.courses, courses) &&
        haveSameCourseIds(state.selectedCourses, syncedSelectedCourses)
      ) {
        return state;
      }

      return {
        courses,
        selectedCourses: syncedSelectedCourses,
        totalSks: calculateTotalSks(syncedSelectedCourses),
      };
    });
  },
  clearSelectedCourses: () => {
    set({
      selectedCourses: [],
      totalSks: 0,
    });
  },
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
      totalSks: calculateTotalSks(nextSelectedCourses),
    });
  },
  removeCourse: (courseId) => {
    const nextSelectedCourses = get().selectedCourses.filter(
      (item) => item.id !== courseId
    );

    set({
      selectedCourses: nextSelectedCourses,
      totalSks: calculateTotalSks(nextSelectedCourses),
    });
  },
  getVisualConflictMap: () => buildVisualConflictMap(get().selectedCourses),
  hasVisualConflict: (courseId) =>
    Boolean(buildVisualConflictMap(get().selectedCourses)[courseId]?.length),
}));
