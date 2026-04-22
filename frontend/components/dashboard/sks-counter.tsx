"use client";

import { animate, motion, useMotionValue, useTransform } from "framer-motion";
import { useEffect } from "react";

import { GlassCard } from "@/components/ui/glass-card";

type SksCounterProps = {
  totalSks: number;
  maxSks?: number;
  selectedCount: number;
  isCoursesLoading?: boolean;
};

export function SksCounter({
  totalSks,
  maxSks = 20,
  selectedCount,
  isCoursesLoading = false,
}: SksCounterProps) {
  const progress = Math.min((totalSks / maxSks) * 100, 100);
  const animatedValue = useMotionValue(totalSks);
  const roundedValue = useTransform(animatedValue, (latest) =>
    Math.round(latest)
  );
  const isOverflow = totalSks > maxSks;

  useEffect(() => {
    const controls = animate(animatedValue, totalSks, {
      duration: 0.45,
      ease: [0.22, 1, 0.36, 1],
    });

    return () => controls.stop();
  }, [animatedValue, totalSks]);

  return (
    <GlassCard className="shrink-0 px-5 py-5 sm:px-6">
      <div className="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.35em] text-zinc-500">
            SKS Monitor
          </p>
          <div className="mt-3 flex items-end gap-3">
            <motion.span className="text-4xl font-semibold tracking-tight text-white">
              {roundedValue}
            </motion.span>
            <span className="pb-1 text-sm text-zinc-500">/ {maxSks} SKS</span>
          </div>
          <p className="mt-3 text-sm leading-6 text-zinc-400">
            Gauge ini memberi batas visual sebelum nanti validasi final dijalankan
            oleh backend Go.
          </p>
          <div className="mt-4 flex items-center gap-3 text-xs">
            <span className="rounded-full border border-white/10 bg-white/[0.04] px-3 py-1 text-zinc-400">
              {selectedCount} course selected
            </span>
            <span
              className={[
                "rounded-full border px-3 py-1",
                isCoursesLoading
                  ? "border-white/10 bg-white/[0.04] text-zinc-400"
                  : "border-emerald-300/20 bg-emerald-300/[0.08] text-emerald-300",
              ].join(" ")}
            >
              {isCoursesLoading ? "Syncing curriculum" : "API connected"}
            </span>
          </div>
        </div>

        <div className="w-full max-w-md">
          <div className="rounded-[1.35rem] border border-white/10 bg-black/20 p-3">
            <div className="flex items-center justify-between text-[11px] uppercase tracking-[0.25em] text-zinc-500">
              <span>Progress</span>
              <span className={isOverflow ? "text-zinc-200" : "text-emerald-300"}>
                {isOverflow ? "Over Limit" : "Healthy"}
              </span>
            </div>

            <div className="mt-4 h-3 overflow-hidden rounded-full bg-white/5">
              <motion.div
                className={[
                  "h-full rounded-full",
                  isOverflow
                    ? "bg-gradient-to-r from-white/35 to-white/70"
                    : "bg-gradient-to-r from-emerald-300/70 via-emerald-200/60 to-white/75",
                ].join(" ")}
                initial={{ width: 0 }}
                animate={{ width: `${progress}%` }}
                transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
              />
            </div>

            <div className="mt-3 flex items-center justify-between text-xs text-zinc-500">
              <span>0 SKS</span>
              <span>{maxSks} SKS</span>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>
  );
}
