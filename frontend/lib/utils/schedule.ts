export interface PositionResult {
  top: number;
  height: number;
}

export interface ScheduleBlock {
  blockId: string;
  day: string;
  startTime: string;
  endTime: string;
}

const START_DAY_MINUTES = 7 * 60;
const END_DAY_MINUTES = 18 * 60;
const TOTAL_DAY_MINUTES = END_DAY_MINUTES - START_DAY_MINUTES; // 660

export const WEEK_DAYS = ["MON", "TUE", "WED", "THU", "FRI", "SAT"];

export function timeToMinutes(time: string): number {
  const [hourStr, minuteStr] = time.split(":");
  const hour = Number(hourStr);
  const minute = Number(minuteStr);

  if (
    Number.isNaN(hour) ||
    Number.isNaN(minute) ||
    hour < 0 ||
    hour > 23 ||
    minute < 0 ||
    minute > 59
  ) {
    throw new Error(`Invalid time format: "${time}"`);
  }

  return hour * 60 + minute;
}

export function normalizeDay(day: string): string {
  const normalized = day.trim().slice(0, 3).toUpperCase();
  if (["MON", "SEN"].includes(normalized)) return "MON";
  if (["TUE", "SEL"].includes(normalized)) return "TUE";
  if (["WED", "RAB"].includes(normalized)) return "WED";
  if (["THU", "KAM"].includes(normalized)) return "THU";
  if (["FRI", "JUM"].includes(normalized)) return "FRI";
  if (["SAT", "SAB"].includes(normalized)) return "SAT";
  return normalized;
}

export function calculatePosition(
  startTime: string,
  endTime: string,
): PositionResult {
  const startMinute = timeToMinutes(startTime);
  const endMinute = timeToMinutes(endTime);

  const clampedStart = Math.max(START_DAY_MINUTES, startMinute);
  const clampedEnd = Math.min(END_DAY_MINUTES, endMinute);

  const top =
    ((clampedStart - START_DAY_MINUTES) / TOTAL_DAY_MINUTES) * 100;
  const height = Math.max(
    0,
    ((clampedEnd - clampedStart) / TOTAL_DAY_MINUTES) * 100,
  );

  return {
    top,
    height,
  };
}

export function hasOverlap(
  firstStart: string,
  firstEnd: string,
  secondStart: string,
  secondEnd: string,
): boolean {
  const aStart = timeToMinutes(firstStart);
  const aEnd = timeToMinutes(firstEnd);
  const bStart = timeToMinutes(secondStart);
  const bEnd = timeToMinutes(secondEnd);

  return aStart < bEnd && bStart < aEnd;
}

export function detectScheduleConflicts(blocks: ScheduleBlock[]): Set<string> {
  const conflictIds = new Set<string>();

  for (let i = 0; i < blocks.length; i += 1) {
    for (let j = i + 1; j < blocks.length; j += 1) {
      const first = blocks[i];
      const second = blocks[j];

      if (normalizeDay(first.day) !== normalizeDay(second.day)) {
        continue;
      }

      if (
        hasOverlap(
          first.startTime,
          first.endTime,
          second.startTime,
          second.endTime,
        )
      ) {
        conflictIds.add(first.blockId);
        conflictIds.add(second.blockId);
      }
    }
  }

  return conflictIds;
}

