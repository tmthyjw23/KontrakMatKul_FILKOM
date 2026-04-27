import LoginPage from "./login/page";

export default function Home() {
    return (
        <LoginPage />
    )
}
"use client";

import { useMemo, useState } from "react";
import axios, { AxiosError } from "axios";
import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import {
  QueryClient,
  QueryClientProvider,
  useMutation,
  useQuery,
} from "@tanstack/react-query";
import { Toaster, toast } from "sonner";
import { useContractStore, type ContractCourse } from "@/lib/store/useContractStore";
import {
  WEEK_DAYS,
  calculatePosition,
  detectScheduleConflicts,
  normalizeDay,
  type ScheduleBlock,
} from "@/lib/utils/schedule";

type ApiCourse = {
  id?: string | number;
  course_id?: string | number;
  code?: string;
  course_code?: string;
  name?: string;
  course_name?: string;
  sks?: number | string;
  credits?: number | string;
  lecturer?: string;
  schedules?: ApiSchedule[];
  schedule?: ApiSchedule[];
};

type ApiSchedule = {
  day?: string;
  day_of_week?: string;
  startTime?: string;
  endTime?: string;
  start_time?: string;
  end_time?: string;
  room?: string;
};

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1",
});

apiClient.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = window.localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

function normalizeCourses(payload: unknown): ContractCourse[] {
  const rows = Array.isArray(payload)
    ? payload
    : (payload as { data?: unknown[]; courses?: unknown[] })?.data ??
      (payload as { data?: unknown[]; courses?: unknown[] })?.courses ??
      [];

  const normalized: ContractCourse[] = [];

  rows.forEach((row) => {
    const course = row as ApiCourse;
    const id = String(course.id ?? course.course_id ?? "");
    if (!id) return;

    const schedules = (course.schedules ?? course.schedule ?? []).map((slot) => ({
      day: String(slot.day ?? slot.day_of_week ?? "MON"),
      startTime: String(slot.startTime ?? slot.start_time ?? "07:00"),
      endTime: String(slot.endTime ?? slot.end_time ?? "08:40"),
      room: slot.room,
    }));

    normalized.push({
      id,
      code: String(course.code ?? course.course_code ?? "MK"),
      name: String(course.name ?? course.course_name ?? "Mata Kuliah"),
      sks: Number(course.sks ?? course.credits ?? 0),
      lecturer: course.lecturer,
      schedules,
    });
  });

  return normalized;
}

function DashboardContent() {
  const router = useRouter();
  const selectedCourses = useContractStore((state) => state.selectedCourses);
  const totalSks = useContractStore((state) => state.totalSks());
  const maxSks = useContractStore((state) => state.maxSks);
  const isSelected = useContractStore((state) => state.isSelected);
  const toggleCourse = useContractStore((state) => state.toggleCourse);
  const clearCourses = useContractStore((state) => state.clearCourses);

  const { data: courses = [], isLoading, isError } = useQuery({
    queryKey: ["courses"],
    queryFn: async () => {
      const response = await apiClient.get("/courses");
      return normalizeCourses(response.data);
    },
  });

  const scheduleBlocks = useMemo<ScheduleBlock[]>(() => {
    return selectedCourses.flatMap((course) =>
      course.schedules.map((slot, index) => ({
        blockId: `${course.id}-${index}`,
        day: slot.day,
        startTime: slot.startTime,
        endTime: slot.endTime,
      })),
    );
  }, [selectedCourses]);

  const conflicts = useMemo(
    () => detectScheduleConflicts(scheduleBlocks),
    [scheduleBlocks],
  );

  const submitMutation = useMutation({
    mutationFn: async () => {
      const courseIds = selectedCourses.map((course) => course.id);
      await apiClient.post("/registrations", { courseIds });
    },
    onSuccess: () => {
      toast.success("Kontrak mata kuliah berhasil dikirim.");
    },
    onError: (error: AxiosError<{ message?: string }>) => {
      toast.error(
        error.response?.data?.message ??
          "Gagal submit kontrak. Periksa koneksi atau login.",
      );
    },
  });

  const onToggleCourse = (course: ContractCourse) => {
    const result = toggleCourse(course);
    if (!result.ok && result.reason === "max_sks") {
      toast.error("Maksimal 24 SKS. Hapus mata kuliah lain terlebih dahulu.");
    }
  };

  const hasConflict = conflicts.size > 0;

  const showProfile = () => {
    const fallbackProfile = { name: "Mahasiswa", email: "-" };
    const profile = (() => {
      const rawUser = window.localStorage.getItem("user");
      if (!rawUser) return fallbackProfile;
      try {
        const parsed = JSON.parse(rawUser) as {
          name?: string;
          fullName?: string;
          email?: string;
        };
        return {
          name: parsed.name ?? parsed.fullName ?? "Mahasiswa",
          email: parsed.email ?? "-",
        };
      } catch {
        return fallbackProfile;
      }
    })();

    toast.info(`${profile.name} • ${profile.email}`, {
      description: "Informasi profil aktif",
    });
  };

  const handleLogout = () => {
    window.localStorage.removeItem("token");
    window.localStorage.removeItem("user");
    clearCourses();
    toast.success("Logout berhasil.");
    router.refresh();
  };

  return (
    <div className="relative min-h-screen overflow-hidden bg-zinc-950 px-4 py-6 text-zinc-100 sm:px-8">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_10%_20%,rgba(255,255,255,0.12),transparent_40%),radial-gradient(circle_at_90%_80%,rgba(255,255,255,0.08),transparent_35%)]" />
      <main className="relative mx-auto grid max-w-7xl grid-cols-1 gap-4 lg:grid-cols-[1.1fr_1.4fr]">
        <section className="rounded-3xl border border-zinc-100/10 bg-zinc-800/25 p-4 shadow-2xl backdrop-blur-xl">
          <div className="mb-4 flex items-start justify-between gap-3">
            <div>
              <p className="text-xs uppercase tracking-[0.16em] text-zinc-400">
                FILKOM UNKLAB
              </p>
              <h1 className="text-xl font-semibold text-white">
                Dashboard Kontrak Mata Kuliah
              </h1>
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={showProfile}
                className="rounded-xl border border-zinc-100/20 bg-zinc-900/70 px-3 py-2 text-xs font-medium text-zinc-100 transition hover:border-zinc-300/40"
              >
                Profile
              </button>
              <button
                type="button"
                onClick={handleLogout}
                className="rounded-xl border border-zinc-100/20 bg-zinc-100 px-3 py-2 text-xs font-semibold text-zinc-900 transition hover:bg-zinc-200"
              >
                Logout
              </button>
              <div className="rounded-2xl border border-zinc-100/10 bg-zinc-950/50 px-3 py-2 text-right">
                <p className="text-xs text-zinc-400">Total SKS</p>
                <p className="text-xl font-bold text-white">
                  {totalSks}/{maxSks}
                </p>
              </div>
            </div>
          </div>

          {isLoading && <p className="text-sm text-zinc-300">Memuat mata kuliah...</p>}
          {isError && (
            <p className="text-sm text-red-300">
              Data mata kuliah gagal dimuat dari backend.
            </p>
          )}

          <div className="grid grid-cols-1 gap-3">
            {courses.map((course, idx) => {
              const active = isSelected(course.id);
              return (
                <motion.button
                  key={course.id}
                  type="button"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: idx * 0.02 }}
                  onClick={() => onToggleCourse(course)}
                  className={`rounded-2xl border p-4 text-left transition ${
                    active
                      ? "border-zinc-200 bg-zinc-200/10"
                      : "border-zinc-100/10 bg-zinc-900/40 hover:border-zinc-400/40"
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <p className="text-sm font-semibold text-white">
                      {course.code} - {course.name}
                    </p>
                    <span className="rounded-full border border-zinc-300/20 px-2 py-1 text-xs text-zinc-200">
                      {course.sks} SKS
                    </span>
                  </div>
                  {course.lecturer ? (
                    <p className="mt-1 text-xs text-zinc-400">{course.lecturer}</p>
                  ) : null}
                </motion.button>
              );
            })}
          </div>

          <button
            type="button"
            disabled={submitMutation.isPending || selectedCourses.length === 0 || hasConflict}
            onClick={() => submitMutation.mutate()}
            className="mt-4 w-full rounded-2xl border border-zinc-100/20 bg-zinc-100 px-4 py-3 text-sm font-semibold text-zinc-900 disabled:cursor-not-allowed disabled:opacity-40"
          >
            {submitMutation.isPending ? "Mengirim..." : "Submit Kontrak"}
          </button>
          {hasConflict ? (
            <p className="mt-2 text-xs text-red-300">
              Terdapat konflik jadwal. Perbaiki sebelum submit.
            </p>
          ) : null}
        </section>

        <section className="rounded-3xl border border-zinc-100/10 bg-zinc-900/30 p-4 shadow-2xl backdrop-blur-xl">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Schedule Grid</h2>
            <p className="text-xs text-zinc-400">07:00 - 18:00</p>
          </div>

          <div className="grid grid-cols-[52px_repeat(6,minmax(80px,1fr))] gap-2">
            <div />
            {WEEK_DAYS.map((day) => (
              <div
                key={day}
                className="rounded-xl border border-zinc-100/10 bg-zinc-800/40 py-2 text-center text-xs font-medium text-zinc-200"
              >
                {day}
              </div>
            ))}

            <div className="relative h-[640px]">
              {Array.from({ length: 12 }, (_, index) => index + 7).map((hour) => (
                <div
                  key={hour}
                  className="absolute -translate-y-1/2 text-[10px] text-zinc-400"
                  style={{ top: `${((hour - 7) / 11) * 100}%` }}
                >
                  {String(hour).padStart(2, "0")}:00
                </div>
              ))}
            </div>

            {WEEK_DAYS.map((day) => (
              <div
                key={day}
                className="relative h-[640px] overflow-hidden rounded-2xl border border-zinc-100/10 bg-zinc-950/60"
              >
                {Array.from({ length: 12 }, (_, index) => (
                  <div
                    key={index}
                    className="absolute left-0 right-0 border-t border-zinc-100/5"
                    style={{ top: `${(index / 11) * 100}%` }}
                  />
                ))}

                {selectedCourses.map((course) =>
                  course.schedules
                    .filter((slot) => normalizeDay(slot.day) === day)
                    .map((slot, index) => {
                      const blockId = `${course.id}-${index}`;
                      const { top, height } = calculatePosition(
                        slot.startTime,
                        slot.endTime,
                      );
                      const conflict = conflicts.has(blockId);

                      return (
                        <motion.div
                          key={blockId}
                          initial={{ opacity: 0, scale: 0.95 }}
                          animate={{ opacity: 1, scale: 1 }}
                          className={`absolute left-1 right-1 rounded-xl border p-2 text-[10px] shadow-lg ${
                            conflict
                              ? "border-red-300/70 bg-red-500/60 text-white"
                              : "border-zinc-200/20 bg-zinc-100/80 text-zinc-900"
                          }`}
                          style={{ top: `${top}%`, height: `${height}%` }}
                        >
                          <p className="font-bold leading-tight">{course.code}</p>
                          <p className="leading-tight">{slot.startTime}</p>
                          <p className="leading-tight">{slot.endTime}</p>
                        </motion.div>
                      );
                    }),
                )}
              </div>
            ))}
          </div>
        </section>
      </main>
      <Toaster richColors theme="dark" position="top-right" />
    </div>
  );
}

export default function Page() {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <DashboardContent />
    </QueryClientProvider>
  );
}
