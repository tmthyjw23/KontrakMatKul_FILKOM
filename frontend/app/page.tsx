"use client";

import { CourseCard } from "@/components/dashboard/course-card";
import { ScheduleGrid } from "@/components/dashboard/schedule-grid";
import { SksCounter } from "@/components/dashboard/sks-counter";
import { GlassCard } from "@/components/ui/glass-card";
import { useContractStore } from "@/lib/store/useContractStore";
import type { ContractState } from "@/lib/store/useContractStore";

export default function Home() {
  const courses = useContractStore((state: ContractState) => state.courses);
  const selectedCourses = useContractStore(
    (state: ContractState) => state.selectedCourses
  );
  const totalSks = useContractStore((state: ContractState) => state.totalSks);
  const maxSks = useContractStore((state: ContractState) => state.maxSks);
  const isConfirmDisabled = totalSks === 0 || totalSks > maxSks;

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
              <div className="space-y-3">
                {courses.map((course) => (
                  <CourseCard key={course.id} course={course} />
                ))}
              </div>
            </div>
          </GlassCard>
        </section>

        <section className="col-span-12 flex min-h-0 flex-col gap-4 xl:col-span-8">
          <SksCounter totalSks={totalSks} maxSks={maxSks} />

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
                  disabled={isConfirmDisabled}
                  className={[
                    "rounded-[1.35rem] border px-5 py-4 text-sm font-medium transition",
                    isConfirmDisabled
                      ? "cursor-not-allowed border-white/10 bg-white/[0.05] text-zinc-500"
                      : "border-emerald-300/30 bg-emerald-300/[0.12] text-emerald-100 hover:border-emerald-200/45 hover:bg-emerald-200/[0.16]",
                  ].join(" ")}
                >
                  Confirm Enrollment
                </button>
              </div>
            </div>
          </GlassCard>
        </section>
      </div>
    </main>
  );
}
