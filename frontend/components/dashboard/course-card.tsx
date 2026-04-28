"use client";

import { motion } from "framer-motion";

import { useContractStore } from "@/lib/store/useContractStore";
import type { ContractState } from "@/lib/store/useContractStore";
import type { Course, Schedule } from "@/src/types/course";

type CourseCardProps = {
  course: Course;
};

function formatSchedules(schedules: Schedule[]) {
  return schedules
    .map(
      (schedule) =>
        `${schedule.day}, ${schedule.start_time} - ${schedule.end_time} • ${schedule.room}`
    )
    .join(" • ");
}

export function CourseCard({ course }: CourseCardProps) {
  const addCourse = useContractStore((state) => state.addCourse);
  const removeCourse = useContractStore((state) => state.removeCourse);
  const selectedCourses = useContractStore(
    (state: ContractState) => state.selectedCourses
  );
  const hasVisualConflict = useContractStore(
    (state) => state.hasVisualConflict
  );

  const isSelected = selectedCourses.some((item) => item.id === course.id);
  const isConflicted = isSelected && hasVisualConflict(course.id);

  const handleToggle = () => {
    if (isSelected) {
      removeCourse(course.id);
    } else {
      addCourse(course.id);
    }
  };

  return (
    <motion.button
      type="button"
      whileHover={{ scale: 1.015 }}
      whileTap={{ scale: 0.995 }}
      onClick={handleToggle}
      className={[
        "w-full rounded-[1.5rem] border px-5 py-4 text-left transition-all duration-300",
        "bg-white/[0.04] backdrop-blur-md",
        "hover:border-white/20 hover:bg-white/[0.07]",
        isConflicted
          ? "border-white/30 bg-white/[0.11] shadow-[0_18px_60px_rgba(0,0,0,0.4)]"
          : isSelected
          ? "border-emerald-300/25 bg-emerald-300/[0.08] shadow-[0_18px_60px_rgba(0,0,0,0.35)]"
          : "border-white/10",
      ].join(" ")}
    >
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-xs font-medium uppercase tracking-[0.3em] text-zinc-400">
            {course.code}
          </p>
          <h3 className="mt-2 text-lg font-semibold text-white">
            {course.name}
          </h3>
        </div>
        <div className="rounded-full border border-white/10 bg-white/[0.05] px-3 py-1 text-xs font-medium text-zinc-200">
          {course.sks} SKS
        </div>
      </div>

      <div className="mt-4 space-y-2 text-sm text-zinc-300">
        <p>
          <span className="text-zinc-500">Lecturer</span>{" "}
          <span className="text-zinc-100">{course.lecturer}</span>
        </p>
        <p className="leading-6 text-zinc-400">
          {formatSchedules(course.schedules)}
        </p>
      </div>

      <div className="mt-4 flex items-center justify-between text-xs">
        <span className="text-zinc-500">
          {isConflicted
            ? "Ada bentrok visual di jadwal"
            : "Klik untuk memetakan ke jadwal mingguan"}
        </span>
        <span
          className={[
            isSelected ? "font-medium" : "text-zinc-500",
            isConflicted
              ? "text-zinc-100"
              : isSelected
              ? "text-emerald-300"
              : "",
          ].join(" ")}
        >
          {isConflicted ? "Conflict" : isSelected ? "Selected" : "Available"}
        </span>
      </div>
    </motion.button>
  );
}
