"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

import { studentApi } from "@/lib/api/admin";
import { useAuthStore } from "@/lib/store/useAuthStore";
import type { Course } from "@/src/types/course";

type EnrollmentSuccess = { count: number };

function extractApiErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data;
    // Backend may return { error: "..." } or { message: "..." }
    const message = data?.message ?? data?.error;
    if (typeof message === "string") return message;
  }
  if (error instanceof Error) return error.message;
  return "Failed to submit enrollment.";
}

async function submitEnrollment(
  selectedCourses: Course[],
  nim: string
): Promise<EnrollmentSuccess> {
  if (!nim) {
    throw new Error("User NIM not available. Please log in again.");
  }

  for (const course of selectedCourses) {
    try {
      await studentApi.registerCourse(nim, course.code);
    } catch (error) {
      const message = extractApiErrorMessage(error);
      throw new Error(
        message.includes(course.code) ? message : `${message} (${course.code})`
      );
    }
  }

  return { count: selectedCourses.length };
}

export function useEnrollment() {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const nim = user?.student_number ?? "";

  return useMutation({
    mutationFn: (selectedCourses: Course[]) =>
      submitEnrollment(selectedCourses, nim),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["courses"] });
    },
  });
}
