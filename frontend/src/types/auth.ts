export type UserRole = "student" | "admin";

export interface User {
  id: string;
  name: string;
  email?: string;
  student_number?: string;
  role: UserRole;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

export interface LoginRequest {
  student_number?: string;
  password: string;
  role: UserRole;
}

// Raw data from backend login endpoint (after Antigravity adds response wrapper)
// Backend: POST /api/v1/auth/login → { token, role, expires_at }
// Antigravity will wrap as: { code, status, message, data: LoginData }
export interface LoginData {
  token: string;
  role: string;       // "Student" or "Admin" (capitalized, from backend)
  expires_at: number; // Unix timestamp
  // Antigravity may add these to avoid an extra profile fetch:
  name?: string;
  nim?: string;
}

export interface LoginResponse {
  code: number;
  status: string;
  message: string;
  data: LoginData;
}

export interface ContractPeriod {
  id: string;
  is_open: boolean;
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at: string;
}

// Raw backend Registration from GET /api/v1/admin/registrations
// Antigravity will enrich this with student_name, course_name via JOIN
export interface BackendRegistration {
  id: number;
  student_nim: string;
  course_code: string;
  status: string; // "registered" | "cancelled" | "pending" | "approved" | "rejected"
  created_at: string;
  // Enriched fields (added by Antigravity via JOIN query):
  student_name?: string;
  student_number?: string;
  course_name?: string;
}

export interface StudentEnrollment {
  id: string;
  student_id: string;
  student_name: string;
  student_number: string;
  course_id: string;
  course_code: string;
  course_name: string;
  status: "pending" | "approved" | "rejected" | "registered" | "cancelled";
  enrolled_at: string;
}
