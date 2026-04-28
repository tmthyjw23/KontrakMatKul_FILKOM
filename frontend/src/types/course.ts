export interface Schedule {
  day: string;
  start_time: string;
  end_time: string;
  room: string;
}

export interface Course {
  id: string;
  code: string;
  name: string;
  sks: number;
  quota?: number; // maps to backend cohort_target
  lecturer: string;
  schedules: Schedule[];
  color?: string;
}

export interface VisualConflictMap {
  [courseId: string]: string[];
}
