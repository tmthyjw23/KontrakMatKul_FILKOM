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

export interface LoginResponse {
  code: number;
  status: string;
  message: string;
  data: {
    token: string;
    user: User;
  };
}

export interface ContractPeriod {
  id: string;
  is_open: boolean;
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at: string;
}

export interface StudentEnrollment {
  id: string;
  student_id: string;
  student_name: string;
  student_number: string;
  course_id: string;
  course_code: string;
  course_name: string;
  status: "pending" | "approved" | "rejected";
  enrolled_at: string;
}
