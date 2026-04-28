"use client";

import { GlassCard } from "@/components/ui/glass-card";
import { useContractStore } from "@/lib/store/useContractStore";
import type { ContractState } from "@/lib/store/useContractStore";
import {
  calculatePosition,
  formatHourLabel,
  SCHEDULE_DURATION_MINUTES,
  SCHEDULE_START_MINUTES,
  normalizeDay,
} from "@/lib/utils/schedule";

import type { Course, Schedule } from "@/src/types/course";

const DAYS = ["Senin", "Selasa", "Rabu", "Kamis", "Jumat"] as const;
const GRID_HEIGHT = 880;
const TIME_MARKERS = Array.from({ length: 12 }, (_, index) => {
  const minutes = SCHEDULE_START_MINUTES + index * 60;
  return {
    label: formatHourLabel(minutes),
    top:
      ((minutes - SCHEDULE_START_MINUTES) / SCHEDULE_DURATION_MINUTES) * 100,
  };
});

export function ScheduleGrid() {
  const selectedCourses = useContractStore(
    (state: ContractState) => state.selectedCourses
  );
  const removeCourse = useContractStore((state) => state.removeCourse);
  const getVisualConflictMap = useContractStore(
    (state) => state.getVisualConflictMap
  );
  const conflictMap = getVisualConflictMap();

  const blocksByDay = DAYS.reduce<
    Record<string, Array<Course & { slot: Schedule }>>
  >((accumulator, day) => {
    accumulator[day] = selectedCourses.flatMap((course) =>
      course.schedules
        .filter((schedule) => normalizeDay(schedule.day) === day)
        .map((slot) => ({ ...course, slot }))
    );

    return accumulator;
  }, {});

  return (
    <GlassCard className="flex h-full min-h-0 flex-col overflow-hidden p-0">
      <div className="border-b border-white/10 bg-black/20 px-5 py-4 sm:px-6">
        <div className="flex items-end justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.35em] text-zinc-500">
              Weekly Planner
            </p>
            <h2 className="mt-2 text-2xl font-semibold text-white">
              Real-time schedule map
            </h2>
          </div>
          <p className="max-w-xs text-right text-sm leading-6 text-zinc-400">
            Setiap blok memakai posisi absolut, jadi jam mulai dan selesai tetap
            akurat sampai level menit.
          </p>
        </div>
      </div>

      <div className="min-h-0 flex-1 overflow-auto">
        <div className="min-w-[920px] p-4 sm:p-5">
          <div className="grid grid-cols-[80px_repeat(5,minmax(0,1fr))] overflow-hidden rounded-[1.5rem] border border-white/10 bg-black/20">
            <div className="border-r border-white/10 bg-white/[0.03]" />

            {DAYS.map((day) => (
              <div
                key={day}
                className="border-l border-white/10 bg-white/[0.03] px-4 py-4 text-center"
              >
                <p className="text-xs uppercase tracking-[0.3em] text-zinc-500">
                  Day
                </p>
                <p className="mt-2 text-sm font-medium text-zinc-100">{day}</p>
              </div>
            ))}

            <div
              className="relative border-r border-white/10 bg-white/[0.02]"
              style={{ height: GRID_HEIGHT }}
            >
              {TIME_MARKERS.map((marker, index) => (
                <div
                  key={marker.label}
                  className="absolute inset-x-0"
                  style={{ top: `${marker.top}%` }}
                >
                  <span
                    className={[
                      "absolute right-3 text-[11px] font-medium text-zinc-500",
                      index === 0 ? "top-1" : "-translate-y-1/2",
                    ].join(" ")}
                  >
                    {marker.label}
                  </span>
                  <div className="absolute inset-x-0 top-0 border-t border-dashed border-white/6" />
                </div>
              ))}
            </div>

            {DAYS.map((day) => (
              <div
                key={day}
                className="relative border-l border-white/10 bg-[linear-gradient(to_bottom,rgba(255,255,255,0.03),rgba(255,255,255,0.015))]"
                style={{ height: GRID_HEIGHT }}
              >
                {TIME_MARKERS.map((marker) => (
                  <div
                    key={`${day}-${marker.label}`}
                    className="absolute inset-x-0 border-t border-dashed border-white/8"
                    style={{ top: `${marker.top}%` }}
                  />
                ))}

                {blocksByDay[day].map((course) => {
                  const { top, height } = calculatePosition(
                    course.slot.start_time,
                    course.slot.end_time
                  );
                  const isConflicted = Boolean(conflictMap[course.id]?.length);

                  return (
                    <button
                      key={`${course.id}-${day}-${course.slot.start_time}`}
                      type="button"
                      onClick={() => removeCourse(course.id)}
                      className={[
                        "absolute left-2 right-2 overflow-hidden rounded-[1.35rem] border p-3 text-left shadow-[0_18px_50px_rgba(0,0,0,0.3)] backdrop-blur-xl transition",
                        isConflicted
                          ? "border-white/25 bg-white/[0.14] hover:border-white/35 hover:bg-white/[0.16]"
                          : "border-emerald-300/20 bg-emerald-300/[0.09] hover:border-emerald-200/35 hover:bg-emerald-200/[0.12]",
                      ].join(" ")}
                      style={{
                        top: `${top}%`,
                        height: `${height}%`,
                        minHeight: 84,
                      }}
                    >
                      <div className="flex h-full flex-col justify-between gap-3">
                        <div>
                          <p className="text-[11px] uppercase tracking-[0.28em] text-zinc-400">
                            {course.code}
                          </p>
                          <h3 className="mt-2 text-sm font-semibold text-white">
                            {course.name}
                          </h3>
                        </div>

                        <div className="space-y-1 text-xs text-zinc-300">
                          <p>
                            {course.slot.start_time} - {course.slot.end_time}
                          </p>
                          <p>{course.slot.room}</p>
                          <p>{course.sks} SKS</p>
                          <p
                            className={
                              isConflicted ? "text-zinc-300" : "text-emerald-300"
                            }
                          >
                            {isConflicted ? "Visual conflict detected" : "Click to remove"}
                          </p>
                        </div>
                      </div>
                    </button>
                  );
                })}
              </div>
            ))}
          </div>

          {selectedCourses.length === 0 ? (
            <div className="mt-4 rounded-[1.25rem] border border-dashed border-white/10 bg-white/[0.03] px-5 py-4 text-sm text-zinc-400">
              Belum ada mata kuliah yang dipilih. Klik course di panel kiri untuk
              melihatnya langsung muncul di grid.
            </div>
          ) : null}
        </div>
      </div>
    </GlassCard>
  );
}
