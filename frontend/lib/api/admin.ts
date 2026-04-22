"use client";

import { apiClient } from "./client";
import type {
  LoginRequest,
  LoginResponse,
  ContractPeriod,
  StudentEnrollment,
} from "@/src/types/auth";
import type { Course } from "@/src/types/course";

export const authApi = {
  login: async (request: LoginRequest) => {
    const response = await apiClient.post<LoginResponse>("/login", {
      student_number: request.student_number,
      password: request.password,
    });
    return response.data.data;
  },
};

export const adminApi = {
  // Contract Period Management
  getContractPeriod: async () => {
    const response = await apiClient.get<{
      code: number;
      data: ContractPeriod;
    }>("/admin/contract-period");
    return response.data.data;
  },

  updateContractPeriod: async (period: Partial<ContractPeriod>) => {
    const response = await apiClient.put<{
      code: number;
      data: ContractPeriod;
    }>("/admin/contract-period", period);
    return response.data.data;
  },

  toggleContractPeriod: async (isOpen: boolean) => {
    return adminApi.updateContractPeriod({ is_open: isOpen });
  },

  // Courses Management
  getCourses: async () => {
    const response = await apiClient.get<{
      code: number;
      data: Course[];
    }>("/admin/courses");
    return response.data.data;
  },

  createCourse: async (course: Partial<Course>) => {
    const response = await apiClient.post<{
      code: number;
      data: Course;
    }>("/admin/courses", course);
    return response.data.data;
  },

  updateCourse: async (id: string, course: Partial<Course>) => {
    const response = await apiClient.put<{
      code: number;
      data: Course;
    }>(`/admin/courses/${id}`, course);
    return response.data.data;
  },

  deleteCourse: async (id: string) => {
    const response = await apiClient.delete<{
      code: number;
    }>(`/admin/courses/${id}`);
    return response.data;
  },

  // Enrollments Monitoring
  getEnrollments: async () => {
    const response = await apiClient.get<{
      code: number;
      data: StudentEnrollment[];
    }>("/admin/enrollments");
    return response.data.data;
  },

  approveEnrollment: async (id: string) => {
    const response = await apiClient.post<{
      code: number;
      data: StudentEnrollment;
    }>(`/admin/enrollments/${id}/approve`);
    return response.data.data;
  },

  rejectEnrollment: async (id: string) => {
    const response = await apiClient.post<{
      code: number;
      data: StudentEnrollment;
    }>(`/admin/enrollments/${id}/reject`);
    return response.data.data;
  },

  // Users Management
  getStudents: async () => {
    const response = await apiClient.get<{
      code: number;
      data: Array<{
        id: string;
        name: string;
        student_number: string;
        email?: string;
      }>;
    }>("/admin/students");
    return response.data.data;
  },

  createStudent: async (student: {
    name: string;
    student_number: string;
    password: string;
    email?: string;
  }) => {
    const response = await apiClient.post<{
      code: number;
      data: { id: string; name: string; student_number: string };
    }>("/admin/students", student);
    return response.data.data;
  },

  resetStudentPassword: async (id: string, newPassword: string) => {
    const response = await apiClient.post<{
      code: number;
      message: string;
    }>(`/admin/students/${id}/reset-password`, { password: newPassword });
    return response.data;
  },

  deleteStudent: async (id: string) => {
    const response = await apiClient.delete<{
      code: number;
    }>(`/admin/students/${id}`);
    return response.data;
  },
};
