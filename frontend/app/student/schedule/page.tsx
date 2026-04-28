"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { useCourses } from "@/lib/hooks/useCourses";
import { useAuth } from "@/lib/store/useAuthStore";
import { studentApi } from "@/lib/api/admin";
import type { StudentEnrollment } from "@/src/types/auth";
import type { Course } from "@/src/types/course";

const DAYS_OF_WEEK = [
  "Senin",
  "Selasa",
  "Rabu",
  "Kamis",
  "Jumat",
  "Sabtu",
];
const TIME_SLOTS = [
  "07:00",
  "08:00",
  "09:00",
  "10:00",
  "11:00",
  "12:00",
  "13:00",
  "14:00",
  "15:00",
  "16:00",
];

function ScheduleContent() {
  const { data: allCourses, isLoading: coursesLoading } = useCourses();
  const { user } = useAuth();

  const [registrations, setRegistrations] = useState<StudentEnrollment[]>([]);
  const [registrationsLoading, setRegistrationsLoading] = useState(true);

  useEffect(() => {
    const nim = user?.student_number;
    if (!nim) {
      setRegistrationsLoading(false);
      return;
    }

    studentApi
      .getMyRegistrations(nim)
      .then(setRegistrations)
      .catch(() => setRegistrations([]))
      .finally(() => setRegistrationsLoading(false));
  }, [user?.student_number]);

  const isLoading = coursesLoading || registrationsLoading;

  // Only show courses with an active/approved registration
  const enrolledCourseCodes = new Set(
    registrations
      .filter(
        (r) => r.status === "approved" || r.status === "registered"
      )
      .map((r) => r.course_code)
  );

  const enrolledCourses: Course[] =
    allCourses?.filter((c) => enrolledCourseCodes.has(c.code)) ?? [];

  // Build a lookup: day → time slot → course for the grid
  const scheduleMap = new Map<string, Course>();
  for (const course of enrolledCourses) {
    for (const slot of course.schedules) {
      scheduleMap.set(`${slot.day}-${slot.start_time.slice(0, 5)}`, course);
    }
  }

  if (isLoading) {
    return (
      <>
        <AppHeader />
        <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-6xl space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="h-20 animate-pulse rounded-[1.5rem] border border-white/10 bg-white/[0.04]"
              />
            ))}
          </div>
        </main>
      </>
    );
  }

  return (
    <>
      <AppHeader />
      <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-6xl space-y-8">
          {/* Header */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
            className="space-y-2"
          >
            <h1 className="text-3xl font-bold text-white">My Schedule</h1>
            <p className="text-zinc-400">
              Your enrolled courses and weekly schedule
            </p>
          </motion.div>

          {/* Enrolled Courses */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
            className="space-y-4"
          >
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold text-white">
                Enrolled Courses
              </h2>
              <span className="text-sm text-zinc-500">
                {enrolledCourses.length} course
                {enrolledCourses.length !== 1 ? "s" : ""}
              </span>
            </div>

            {enrolledCourses.length > 0 ? (
              <div className="grid gap-4 sm:grid-cols-2">
                {enrolledCourses.map((course) => (
                  <GlassCard
                    key={course.id}
                    className={`p-6 bg-gradient-to-br ${course.color ?? "from-emerald-300/18 via-white/8 to-white/5"}`}
                  >
                    <div className="space-y-3">
                      <div>
                        <p className="text-xs uppercase tracking-widest text-zinc-400">
                          {course.code}
                        </p>
                        <h3 className="mt-1 text-lg font-semibold text-white">
                          {course.name}
                        </h3>
                      </div>

                      <div className="grid grid-cols-2 gap-3 text-sm">
                        <div>
                          <p className="text-xs text-zinc-500">SKS</p>
                          <p className="font-semibold text-zinc-100">
                            {course.sks}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-zinc-500">Lecturer</p>
                          <p className="font-semibold text-zinc-100 truncate">
                            {course.lecturer}
                          </p>
                        </div>
                      </div>

                      {course.schedules.length > 0 && (
                        <div className="border-t border-white/10 pt-3">
                          <p className="text-xs text-zinc-500 mb-2">Schedule</p>
                          <div className="space-y-1">
                            {course.schedules.map((schedule, sidx) => (
                              <div
                                key={sidx}
                                className="text-sm text-zinc-300 flex items-center gap-2"
                              >
                                <span className="font-medium">
                                  {schedule.day}
                                </span>
                                <span className="text-zinc-600">
                                  {schedule.start_time}–{schedule.end_time}
                                </span>
                                <span className="text-xs text-zinc-600">
                                  {schedule.room}
                                </span>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  </GlassCard>
                ))}
              </div>
            ) : (
              <GlassCard className="p-8 text-center">
                <p className="text-zinc-400">
                  No approved courses yet. Go to the Dashboard to register.
                </p>
              </GlassCard>
            )}
          </motion.div>

          {/* Weekly Schedule Grid */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.2 }}
            className="space-y-4"
          >
            <h2 className="text-xl font-semibold text-white">
              Weekly Schedule
            </h2>

            <GlassCard className="overflow-hidden">
              <div className="overflow-x-auto">
                <div
                  className="inline-grid min-w-full"
                  style={{
                    gridTemplateColumns: `5rem repeat(${DAYS_OF_WEEK.length}, minmax(8rem, 1fr))`,
                  }}
                >
                  {/* Header row */}
                  <div className="px-3 py-3 border-b border-white/10 bg-white/[0.02]" />
                  {DAYS_OF_WEEK.map((day) => (
                    <div
                      key={day}
                      className="px-3 py-3 text-xs font-semibold text-zinc-400 uppercase tracking-widest border-l border-b border-white/10 bg-white/[0.02] text-center"
                    >
                      {day.slice(0, 3)}
                    </div>
                  ))}

                  {/* Time slot rows */}
                  {TIME_SLOTS.map((time) => (
                    <>
                      <div
                        key={`time-${time}`}
                        className="px-3 py-4 text-xs font-semibold text-zinc-500 flex items-start border-b border-white/10"
                      >
                        {time}
                      </div>
                      {DAYS_OF_WEEK.map((day) => {
                        const course = scheduleMap.get(`${day}-${time}`);
                        return (
                          <div
                            key={`${day}-${time}`}
                            className={[
                              "px-2 py-4 border-l border-b border-white/10 text-xs",
                              course
                                ? "bg-emerald-500/10 text-emerald-200"
                                : "",
                            ].join(" ")}
                          >
                            {course && (
                              <span className="font-medium leading-tight block">
                                {course.code}
                              </span>
                            )}
                          </div>
                        );
                      })}
                    </>
                  ))}
                </div>
              </div>
            </GlassCard>

            <p className="text-xs text-zinc-500">
              Only approved/registered courses appear in the schedule grid.
            </p>
          </motion.div>

          {/* All Registrations (pending / rejected too) */}
          {registrations.length > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: 0.3 }}
              className="space-y-4"
            >
              <h2 className="text-xl font-semibold text-white">
                Registration History
              </h2>
              <GlassCard className="divide-y divide-white/10 overflow-hidden">
                {registrations.map((reg) => (
                  <div
                    key={reg.id}
                    className="flex items-center justify-between px-6 py-4"
                  >
                    <div>
                      <p className="font-medium text-zinc-100">
                        {reg.course_name}
                      </p>
                      <p className="text-sm text-zinc-500">{reg.course_code}</p>
                    </div>
                    <span
                      className={[
                        "px-3 py-1 rounded-full text-xs font-medium",
                        reg.status === "approved" || reg.status === "registered"
                          ? "bg-emerald-500/20 text-emerald-200"
                          : reg.status === "pending"
                            ? "bg-orange-500/20 text-orange-200"
                            : "bg-red-500/20 text-red-200",
                      ].join(" ")}
                    >
                      {reg.status.charAt(0).toUpperCase() + reg.status.slice(1)}
                    </span>
                  </div>
                ))}
              </GlassCard>
            </motion.div>
          )}
        </div>
      </main>
    </>
  );
}

export default function SchedulePage() {
  return (
    <ProtectedRoute requiredRole="student">
      <div className="flex min-h-screen flex-col">
        <ScheduleContent />
      </div>
    </ProtectedRoute>
  );
}
