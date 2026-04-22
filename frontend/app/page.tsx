"use client";

import { useEffect } from "react";
import { toast } from "sonner";

import { CourseCard } from "@/components/dashboard/course-card";
import { ScheduleGrid } from "@/components/dashboard/schedule-grid";
import { SksCounter } from "@/components/dashboard/sks-counter";
import { GlassCard } from "@/components/ui/glass-card";
import { useCourses } from "@/lib/hooks/useCourses";
import { useEnrollment } from "@/lib/hooks/useEnrollment";
import { useContractStore } from "@/lib/store/useContractStore";
import type { ContractState } from "@/lib/store/useContractStore";

export default function Home() {
  const courses = useContractStore((state: ContractState) => state.courses);
  const setCourses = useContractStore((state: ContractState) => state.setCourses);
  const clearSelectedCourses = useContractStore(
    (state: ContractState) => state.clearSelectedCourses
  );
  const selectedCourses = useContractStore(
    (state: ContractState) => state.selectedCourses
  );
  const totalSks = useContractStore((state: ContractState) => state.totalSks);
  const maxSks = useContractStore((state: ContractState) => state.maxSks);
  const {
    data: fetchedCourses,
    isLoading: isCoursesLoading,
    isError: isCoursesError,
    error: coursesError,
    refetch: refetchCourses,
  } = useCourses();
  const enrollmentMutation = useEnrollment();

  const isConfirmDisabled =
    totalSks === 0 ||
    totalSks > maxSks ||
    enrollmentMutation.isPending ||
    isCoursesLoading;

  useEffect(() => {
    if (!fetchedCourses) {
      return;
    }

    setCourses(fetchedCourses);
  }, [fetchedCourses, setCourses]);

  async function handleConfirmEnrollment() {
    if (selectedCourses.length === 0) {
      return;
    }

    try {
      await enrollmentMutation.mutateAsync(selectedCourses);
      clearSelectedCourses();
      toast.success("Enrollment berhasil disimpan.");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Enrollment gagal disimpan.";
      toast.error(message);
    }
  }

  return (
    <main className="h-screen overflow-hidden px-4 py-4 sm:px-6 sm:py-6">
      <div className="bento-grid h-full min-h-0">
        <section className="col-span-12 min-h-0 xl:col-span-4">
          <GlassCard className="flex h-full min-h-0 flex-col overflow-hidden p-0">
            <div className="border-b border-white/10 bg-black/20 px-5 py-5 sm:px-6">
              <p className="text-xs uppercase tracking-[0.35em] text-zinc-500">
                Curriculum
              </p>
              <h1 className="mt-3 text-3xl font-semibold tracking-tight text-white">
                FILKOM course contract
              </h1>
              <p className="mt-3 max-w-xl text-sm leading-6 text-zinc-400">
                Panel kiri menampilkan mata kuliah yang tersedia. Setiap klik
                langsung mengirim course ke planner mingguan di sebelah kanan.
              </p>
            </div>

            <div className="min-h-0 flex-1 overflow-y-auto px-4 py-4 sm:px-5">
              {isCoursesLoading ? (
                <div className="space-y-3">
                  {Array.from({ length: 4 }).map((_, index) => (
                    <div
                      key={`course-skeleton-${index}`}
                      className="h-40 animate-pulse rounded-[1.5rem] border border-white/10 bg-white/[0.04]"
                    />
                  ))}
                </div>
              ) : isCoursesError ? (
                <div className="rounded-[1.5rem] border border-white/10 bg-white/[0.03] p-5">
                  <p className="text-sm font-medium text-zinc-100">
                    Gagal memuat curriculum.
                  </p>
                  <p className="mt-2 text-sm leading-6 text-zinc-400">
                    {coursesError instanceof Error
                      ? coursesError.message
                      : "Terjadi kesalahan saat mengambil data course."}
                  </p>
                  <button
                    type="button"
                    onClick={() => void refetchCourses()}
                    className="mt-4 rounded-[1rem] border border-white/10 bg-white/[0.05] px-4 py-2 text-sm text-zinc-100 transition hover:border-white/20 hover:bg-white/[0.08]"
                  >
                    Retry
                  </button>
                </div>
              ) : (
                <div className="space-y-3">
                  {courses.map((course) => (
                    <CourseCard key={course.id} course={course} />
                  ))}
                </div>
              )}
            </div>
          </GlassCard>
        </section>

        <section className="col-span-12 flex min-h-0 flex-col gap-4 xl:col-span-8">
          <SksCounter
            totalSks={totalSks}
            maxSks={maxSks}
            selectedCount={selectedCourses.length}
            isCoursesLoading={isCoursesLoading}
          />

          <div className="min-h-0 flex-1">
            <ScheduleGrid />
          </div>

          <GlassCard className="shrink-0 px-5 py-5 sm:px-6">
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.35em] text-zinc-500">
                  Enrollment Action
                </p>
                <p className="mt-2 text-sm leading-6 text-zinc-400">
                  Tombol aktif hanya jika total SKS berada di rentang valid dan
                  ada mata kuliah yang dipilih.
                </p>
              </div>

              <div className="flex items-center gap-3">
                <div className="rounded-[1.25rem] border border-white/10 bg-white/[0.05] px-4 py-3">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-zinc-500">
                    Courses
                  </p>
                  <p className="mt-2 text-lg font-semibold text-zinc-100">
                    {selectedCourses.length}
                  </p>
                </div>

                <button
                  type="button"
                  onClick={() => void handleConfirmEnrollment()}
                  disabled={isConfirmDisabled}
                  className={[
                    "rounded-[1.35rem] border px-5 py-4 text-sm font-medium transition",
                    isConfirmDisabled
                      ? "cursor-not-allowed border-white/10 bg-white/[0.05] text-zinc-500"
                      : "border-emerald-300/30 bg-emerald-300/[0.12] text-emerald-100 hover:border-emerald-200/45 hover:bg-emerald-200/[0.16]",
                  ].join(" ")}
                >
                  {enrollmentMutation.isPending
                    ? "Submitting..."
                    : "Confirm Enrollment"}
                </button>
              </div>
            </div>
          </GlassCard>
        </section>
      </div>
    </main>
  );
}
