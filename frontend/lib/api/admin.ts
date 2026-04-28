"use client";

import { apiClient } from "./client";
import type {
  ContractPeriod,
  StudentEnrollment,
  BackendRegistration,
} from "@/src/types/auth";
import type { Course } from "@/src/types/course";

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

/**
 * Extracts `.data` from a wrapped { code, status, message, data } envelope,
 * or returns the value directly if it is already the payload (pre-wrapper state).
 */
function extractData<T>(responseData: unknown): T {
  if (
    responseData !== null &&
    typeof responseData === "object" &&
    "code" in responseData &&
    "data" in responseData
  ) {
    return (responseData as { data: T }).data;
  }
  return responseData as T;
}

/**
 * Decode JWT payload to extract custom claims without verifying the signature.
 * Used only for reading nim/role immediately after login.
 */
function decodeJwtClaims(token: string): { nim?: string; role?: string } {
  try {
    const part = token.split(".")[1];
    if (!part) return {};
    const padded = part.replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(atob(padded)) as { nim?: string; role?: string };
  } catch {
    return {};
  }
}

function normalizeRole(role: string | undefined): "student" | "admin" {
  return (role ?? "").toLowerCase() === "admin" ? "admin" : "student";
}

// ---------------------------------------------------------------------------
// Course normalization
// ---------------------------------------------------------------------------

const COURSE_COLORS = [
  "from-emerald-300/18 via-white/8 to-white/5",
  "from-white/14 via-white/7 to-emerald-200/10",
  "from-white/10 via-zinc-200/10 to-emerald-300/12",
  "from-zinc-200/10 via-white/6 to-emerald-100/10",
];

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function normalizeCourse(bc: any, index: number): Course {
  return {
    id: bc.id != null ? String(bc.id) : bc.code,
    code: bc.code ?? "",
    name: bc.name ?? "",
    sks: Number(bc.sks ?? bc.credits ?? 0),
    quota: Number(bc.quota ?? bc.cohort_target ?? 0),
    lecturer: bc.lecturer ?? bc.lecturer_name ?? "TBD",
    schedules: (bc.schedules ?? []).map((s: any) => ({
      day: s.day ?? s.day_of_week ?? "TBA",
      start_time: String(s.start_time ?? "00:00").slice(0, 5),
      end_time: String(s.end_time ?? "00:00").slice(0, 5),
      room: s.room ?? "TBA",
    })),
    color: bc.color ?? COURSE_COLORS[index % COURSE_COLORS.length],
  };
}

// ---------------------------------------------------------------------------
// Registration normalization
// ---------------------------------------------------------------------------

function normalizeRegistration(reg: BackendRegistration): StudentEnrollment {
  const validStatuses = [
    "pending",
    "approved",
    "rejected",
    "registered",
    "cancelled",
  ] as const;
  type ValidStatus = (typeof validStatuses)[number];
  const status: ValidStatus = validStatuses.includes(
    reg.status as ValidStatus
  )
    ? (reg.status as ValidStatus)
    : "pending";

  return {
    id: String(reg.id),
    student_id: reg.student_nim,
    student_name: reg.student_name ?? reg.student_nim,
    student_number: reg.student_number ?? reg.student_nim,
    course_id: reg.course_code,
    course_code: reg.course_code,
    course_name: reg.course_name ?? reg.course_code,
    status,
    enrolled_at: reg.created_at,
  };
}

// ---------------------------------------------------------------------------
// Helpers to map frontend schedule shape → backend shape (handles both
// normalized Schedule { day, start_time } and raw form { dayOfWeek, startTime })
// ---------------------------------------------------------------------------

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function toBackendSchedule(s: any) {
  return {
    day_of_week: s.day ?? s.dayOfWeek ?? "Monday",
    start_time: s.start_time ?? s.startTime ?? "08:00",
    end_time: s.end_time ?? s.endTime ?? "09:00",
    room: s.room ?? "",
  };
}

// ---------------------------------------------------------------------------
// AUTH API  — POST /api/v1/auth/login
// ---------------------------------------------------------------------------

export const authApi = {
  login: async (request: {
    student_number?: string;
    password: string;
    role: string;
  }) => {
    const res = await apiClient.post<unknown>("/auth/login", {
      nim: request.student_number, // backend field name is `nim`
      password: request.password,
    });

    const loginData = extractData<{
      token: string;
      role?: string;
      expires_at?: number;
      name?: string;
      nim?: string;
    }>(res.data);

    const token = loginData.token;
    const claims = decodeJwtClaims(token);
    const nim = loginData.nim ?? claims.nim ?? request.student_number ?? "";
    const role = normalizeRole(loginData.role ?? claims.role);
    const name = loginData.name ?? nim;

    return {
      token,
      user: {
        id: nim,
        name,
        student_number: nim,
        role,
      },
    };
  },

  // GET /api/v1/student/profile/{nim}
  getProfile: async (nim: string) => {
    const res = await apiClient.get<unknown>(`/student/profile/${nim}`);
    return extractData<{
      nim: string;
      name: string;
      faculty?: string;
      study_program?: string;
      cohort_year?: number;
      role?: string;
    }>(res.data);
  },
};

// ---------------------------------------------------------------------------
// ADMIN API
// ---------------------------------------------------------------------------

export const adminApi = {
  // ── Contract Period ──────────────────────────────────────────────────────
  getContractPeriod: async (): Promise<ContractPeriod> => {
    const res = await apiClient.get<unknown>("/admin/contract-period");
    return extractData<ContractPeriod>(res.data);
  },

  updateContractPeriod: async (
    period: Partial<ContractPeriod>
  ): Promise<ContractPeriod> => {
    const res = await apiClient.put<unknown>("/admin/contract-period", period);
    return extractData<ContractPeriod>(res.data);
  },

  toggleContractPeriod: (isOpen: boolean) =>
    adminApi.updateContractPeriod({ is_open: isOpen }),

  // ── Courses ──────────────────────────────────────────────────────────────
  // Uses public GET /courses (no admin-only list endpoint exists yet).
  // Antigravity can add GET /admin/courses later for protected access.
  getCourses: async (): Promise<Course[]> => {
    const res = await apiClient.get<unknown>("/courses");
    const raw = extractData<unknown>(res.data);
    const list: unknown[] = Array.isArray(raw)
      ? raw
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      : ((raw as any)?.courses ?? []);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (list as any[]).map(normalizeCourse);
  },

  createCourse: async (
    course: Partial<Course> & { quota?: number }
  ): Promise<unknown> => {
    const res = await apiClient.post<unknown>("/admin/courses", {
      code: course.code,
      name: course.name,
      credits: course.sks, // sks → credits
      cohort_target: course.quota,
      lecturer_name: course.lecturer, // lecturer → lecturer_name
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      schedules: ((course as any).schedules ?? []).map(toBackendSchedule),
    });
    return extractData<unknown>(res.data);
  },

  // Backend uses course `code` as the path param, not a UUID id.
  // After normalizeCourse, course.id === course.code, so callers pass course.id correctly.
  updateCourse: async (
    code: string,
    course: Partial<Course> & { quota?: number }
  ): Promise<unknown> => {
    const res = await apiClient.put<unknown>(`/admin/courses/${code}`, {
      name: course.name,
      credits: course.sks,
      cohort_target: course.quota,
      lecturer_name: course.lecturer,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      schedules: ((course as any).schedules ?? []).map(toBackendSchedule),
    });
    return extractData<unknown>(res.data);
  },

  deleteCourse: async (code: string): Promise<unknown> => {
    const res = await apiClient.delete<unknown>(`/admin/courses/${code}`);
    return extractData<unknown>(res.data);
  },

  // ── Registrations (= Enrollments) ────────────────────────────────────────
  getEnrollments: async (): Promise<StudentEnrollment[]> => {
    const res = await apiClient.get<unknown>("/admin/registrations");
    const raw = extractData<unknown>(res.data);
    const list: BackendRegistration[] = Array.isArray(raw) ? raw : [];
    return list.map(normalizeRegistration);
  },

  approveEnrollment: async (id: string): Promise<unknown> => {
    const res = await apiClient.post<unknown>(
      `/admin/registrations/${id}/approve`
    );
    return extractData<unknown>(res.data);
  },

  rejectEnrollment: async (id: string): Promise<unknown> => {
    const res = await apiClient.post<unknown>(
      `/admin/registrations/${id}/reject`
    );
    return extractData<unknown>(res.data);
  },

  // ── Students ─────────────────────────────────────────────────────────────
  getStudents: async (): Promise<
    Array<{ id: string; name: string; student_number: string; email?: string }>
  > => {
    const res = await apiClient.get<unknown>("/admin/students");
    const raw = extractData<unknown>(res.data);
    const list: unknown[] = Array.isArray(raw) ? raw : [];
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (list as any[]).map((s) => ({
      id: s.nim ?? s.id ?? "",
      name: s.name ?? "",
      student_number: s.nim ?? s.student_number ?? "",
      email: s.email,
    }));
  },

  createStudent: async (student: {
    name: string;
    student_number: string;
    password: string;
    email?: string;
  }): Promise<unknown> => {
    const res = await apiClient.post<unknown>("/admin/students", {
      nim: student.student_number, // student_number → nim
      name: student.name,
      password: student.password,
      email: student.email,
    });
    return extractData<unknown>(res.data);
  },

  resetStudentPassword: async (
    nim: string,
    newPassword: string
  ): Promise<unknown> => {
    const res = await apiClient.post<unknown>(
      `/admin/students/${nim}/reset-password`,
      { password: newPassword }
    );
    return extractData<unknown>(res.data);
  },

  deleteStudent: async (nim: string): Promise<unknown> => {
    const res = await apiClient.delete<unknown>(`/admin/students/${nim}`);
    return extractData<unknown>(res.data);
  },
};

// ---------------------------------------------------------------------------
// STUDENT API
// ---------------------------------------------------------------------------

export const studentApi = {
  // Contract period — students need read access to know if registration is open.
  // Reuses the admin endpoint until Antigravity adds a public /contract-period route.
  getContractPeriod: async (): Promise<ContractPeriod> => {
    const res = await apiClient.get<unknown>("/admin/contract-period");
    return extractData<ContractPeriod>(res.data);
  },

  // GET /api/v1/student/registrations/{nim}
  getMyRegistrations: async (nim: string): Promise<StudentEnrollment[]> => {
    const res = await apiClient.get<unknown>(`/student/registrations/${nim}`);
    const raw = extractData<unknown>(res.data);
    const list: BackendRegistration[] = Array.isArray(raw) ? raw : [];
    return list.map(normalizeRegistration);
  },

  // POST /api/v1/student/courses/register  →  { nim, course_code }
  registerCourse: async (nim: string, courseCode: string): Promise<unknown> => {
    const res = await apiClient.post<unknown>("/student/courses/register", {
      nim,
      course_code: courseCode,
    });
    return extractData<unknown>(res.data);
  },

  // DELETE /api/v1/student/registrations/{id}
  cancelRegistration: async (id: string | number): Promise<unknown> => {
    const res = await apiClient.delete<unknown>(
      `/student/registrations/${id}`
    );
    return extractData<unknown>(res.data);
  },
};
