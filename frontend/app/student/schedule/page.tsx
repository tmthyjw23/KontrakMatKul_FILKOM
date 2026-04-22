"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { useCourses } from "@/lib/hooks/useCourses";

const DAYS_OF_WEEK = [
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];
const TIME_SLOTS = [
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
  const { data: courses, isLoading } = useCourses();
  const [takenCourses] = useState<string[]>([
    // Mock data - in real app, fetch from backend
  ]);

  // For demonstration, use first few courses as "taken"
  const enrolledCourses = courses?.slice(0, 2) ?? [];

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
              View your course schedule and enrolled courses
            </p>
          </motion.div>

          {/* Enrolled Courses Section */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
            className="space-y-4"
          >
            <h2 className="text-xl font-semibold text-white">
              Enrolled Courses
            </h2>
            <div className="grid gap-4 sm:grid-cols-2">
              {enrolledCourses.length > 0 ? (
                enrolledCourses.map((course, idx) => (
                  <GlassCard
                    key={course.id}
                    className={`p-6 bg-gradient-to-br ${course.color || "from-emerald-300/18 via-white/8 to-white/5"}`}
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
                          <p className="font-semibold text-zinc-100">
                            {course.lecturer}
                          </p>
                        </div>
                      </div>

                      {/* Schedule Info */}
                      <div className="border-t border-white/10 pt-3">
                        <p className="text-xs text-zinc-500 mb-2">Schedule</p>
                        <div className="space-y-1">
                          {course.schedules.map((schedule, sidx) => (
                            <div
                              key={sidx}
                              className="text-sm text-zinc-300 flex items-center gap-2"
                            >
                              <span className="font-medium">{schedule.day}</span>
                              <span className="text-zinc-600">
                                {schedule.start_time}-{schedule.end_time}
                              </span>
                              <span className="text-xs text-zinc-600">
                                {schedule.room}
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </GlassCard>
                ))
              ) : (
                <GlassCard className="col-span-full p-8 text-center">
                  <p className="text-zinc-400">No courses enrolled yet</p>
                </GlassCard>
              )}
            </div>
          </motion.div>

          {/* Weekly Schedule Grid */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.2 }}
            className="space-y-4"
          >
            <h2 className="text-xl font-semibold text-white">Weekly Schedule</h2>

            <GlassCard className="overflow-hidden">
              <div className="overflow-x-auto">
                <div className="inline-block min-w-full">
                  {/* Header */}
                  <div className="grid border-b border-white/10 bg-white/[0.02]">
                    <div className="w-20 px-4 py-3 text-xs font-semibold text-zinc-400 uppercase tracking-widest" />
                    {DAYS_OF_WEEK.map((day) => (
                      <div
                        key={day}
                        className="w-32 px-4 py-3 text-xs font-semibold text-zinc-400 uppercase tracking-widest border-l border-white/10 text-center"
                      >
                        {day}
                      </div>
                    ))}
                  </div>

                  {/* Time slots */}
                  {TIME_SLOTS.map((time) => (
                    <div key={time} className="grid border-b border-white/10">
                      <div className="w-20 px-4 py-3 text-xs font-semibold text-zinc-500 flex items-center">
                        {time}
                      </div>
                      {DAYS_OF_WEEK.map((day) => (
                        <div
                          key={`${day}-${time}`}
                          className="w-32 px-4 py-8 border-l border-white/10"
                        />
                      ))}
                    </div>
                  ))}
                </div>
              </div>
            </GlassCard>

            <p className="text-xs text-zinc-500">
              💡 Tip: This schedule view will automatically populate with your
              enrolled courses.
            </p>
          </motion.div>

          {/* Course History */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.3 }}
            className="space-y-4"
          >
            <h2 className="text-xl font-semibold text-white">
              Course History
            </h2>
            <GlassCard className="p-8">
              <div className="space-y-4">
                <p className="text-sm text-zinc-400 mb-4">
                  Completed courses will appear here
                </p>

                <div className="border-t border-white/10 pt-4">
                  <p className="text-xs text-zinc-600">
                    No course history yet. Start contracting courses!
                  </p>
                </div>
              </div>
            </GlassCard>
          </motion.div>
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
