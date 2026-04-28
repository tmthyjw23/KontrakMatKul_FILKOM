"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

import { apiClient } from "@/lib/api/client";
import type { Course } from "@/src/types/course";

type EnrollmentRequest = {
  course_id: number;
};

type ApiResponse<T> = {
  code: number;
  status: string;
  message: string;
  data: T;
};

type EnrollmentResult = {
  enrollment: {
    id: number;
    user_id: number;
    course_id: number;
    enrolled_at: string;
  };
};

type EnrollmentSuccess = {
  results: EnrollmentResult[];
};

function toEnrollmentPayload(course: Course): EnrollmentRequest {
  const numericCourseId = Number(course.id);

  if (!Number.isFinite(numericCourseId)) {
    throw new Error(`Invalid course id for enrollment: ${course.code}`);
  }

  return {
    course_id: numericCourseId,
  };
}

function extractApiErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message;
    const status = error.response?.status;

    if ((status === 400 || status === 409) && typeof message === "string") {
      return message;
    }

    if (typeof message === "string") {
      return message;
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "Failed to submit enrollment.";
}

async function submitEnrollment(selectedCourses: Course[]): Promise<EnrollmentSuccess> {
  const results: EnrollmentResult[] = [];

  for (const course of selectedCourses) {
    try {
      const payload = toEnrollmentPayload(course);
      const response = await apiClient.post<ApiResponse<EnrollmentResult>>(
        "/enrollments",
        payload
      );

      results.push(response.data.data);
    } catch (error) {
      const message = extractApiErrorMessage(error);
      throw new Error(
        message.includes(course.code) ? message : `${message} (${course.code})`
      );
    }
  }

  return { results };
}

export function useEnrollment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: submitEnrollment,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["courses"] });
    },
  });
}
