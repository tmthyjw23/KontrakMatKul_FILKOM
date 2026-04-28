export function normalizeDay(day: string): string {
  const mapping: Record<string, string> = {
    "monday": "Senin",
    "senin": "Senin",
    "mon": "Senin",
    "tuesday": "Selasa",
    "selasa": "Selasa",
    "tue": "Selasa",
    "wednesday": "Rabu",
    "rabu": "Rabu",
    "wed": "Rabu",
    "thursday": "Kamis",
    "kamis": "Kamis",
    "thu": "Kamis",
    "friday": "Jumat",
    "jumat": "Jumat",
    "fri": "Jumat",
    "saturday": "Sabtu",
    "sabtu": "Sabtu",
    "sat": "Sabtu",
    "sunday": "Minggu",
    "minggu": "Minggu",
    "sun": "Minggu",
  };

  const normalized = mapping[day.toLowerCase().trim()];
  return normalized ?? day;
}

export const SCHEDULE_START_TIME = "07:00";

export const SCHEDULE_END_TIME = "18:00";
export const SCHEDULE_START_MINUTES = 7 * 60;
export const SCHEDULE_END_MINUTES = 18 * 60;
export const SCHEDULE_DURATION_MINUTES =
  SCHEDULE_END_MINUTES - SCHEDULE_START_MINUTES;

export type SchedulePosition = {
  top: number;
  height: number;
};

export const SAFE_SCHEDULE_POSITION: SchedulePosition = {
  top: 0,
  height: 0,
};

export function timeToMinutes(time: string): number | null {
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

export function calculatePosition(
  startTime: string,
  endTime: string
): SchedulePosition {
  const startMinutes = timeToMinutes(startTime);
  const endMinutes = timeToMinutes(endTime);

  if (startMinutes === null || endMinutes === null) {
    return SAFE_SCHEDULE_POSITION;
  }

  const clampedStart = Math.min(
    Math.max(startMinutes, SCHEDULE_START_MINUTES),
    SCHEDULE_END_MINUTES
  );
  const clampedEnd = Math.min(
    Math.max(endMinutes, SCHEDULE_START_MINUTES),
    SCHEDULE_END_MINUTES
  );

  if (clampedEnd <= clampedStart) {
    return SAFE_SCHEDULE_POSITION;
  }

  const top =
    ((clampedStart - SCHEDULE_START_MINUTES) / SCHEDULE_DURATION_MINUTES) * 100;
  const height =
    ((clampedEnd - clampedStart) / SCHEDULE_DURATION_MINUTES) * 100;

  return { top, height };
}

export function formatHourLabel(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  return `${hours.toString().padStart(2, "0")}:00`;
}
