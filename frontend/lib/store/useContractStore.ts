"use client";

import { create } from "zustand";

export interface CourseSchedule {
  day: string;
  startTime: string;
  endTime: string;
  room?: string;
}

export interface ContractCourse {
  id: string;
  code: string;
  name: string;
  sks: number;
  lecturer?: string;
  schedules: CourseSchedule[];
}

interface AddCourseResult {
  ok: boolean;
  reason?: "duplicate" | "max_sks";
}

interface ContractStore {
  selectedCourses: ContractCourse[];
  maxSks: number;
  totalSks: () => number;
  isSelected: (courseId: string) => boolean;
  addCourse: (course: ContractCourse) => AddCourseResult;
  removeCourse: (courseId: string) => void;
  toggleCourse: (course: ContractCourse) => AddCourseResult;
  clearCourses: () => void;
}

export const useContractStore = create<ContractStore>((set, get) => ({
  selectedCourses: [],
  maxSks: 24,
  totalSks: () =>
    get().selectedCourses.reduce((acc, course) => acc + course.sks, 0),
  isSelected: (courseId) =>
    get().selectedCourses.some((course) => course.id === courseId),
  addCourse: (course) => {
    const { selectedCourses, maxSks } = get();
    const alreadySelected = selectedCourses.some(
      (selected) => selected.id === course.id,
    );

    if (alreadySelected) {
      return { ok: false, reason: "duplicate" };
    }

    const currentSks = selectedCourses.reduce(
      (acc, selected) => acc + selected.sks,
      0,
    );
    if (currentSks + course.sks > maxSks) {
      return { ok: false, reason: "max_sks" };
    }

    set({ selectedCourses: [...selectedCourses, course] });
    return { ok: true };
  },
  removeCourse: (courseId) => {
    set((state) => ({
      selectedCourses: state.selectedCourses.filter(
        (course) => course.id !== courseId,
      ),
    }));
  },
  toggleCourse: (course) => {
    const { isSelected, removeCourse, addCourse } = get();
    if (isSelected(course.id)) {
      removeCourse(course.id);
      return { ok: true };
    }
    return addCourse(course);
  },
  clearCourses: () => set({ selectedCourses: [] }),
}));

