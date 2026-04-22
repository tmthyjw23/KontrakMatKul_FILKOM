"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { Input, Toggle } from "@/components/ui/form";
import { adminApi } from "@/lib/api/admin";
import type { ContractPeriod } from "@/src/types/auth";

function ContractPeriodContent() {
  const [period, setPeriod] = useState<ContractPeriod | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [formData, setFormData] = useState({
    is_open: false,
    start_date: "",
    end_date: "",
  });

  useEffect(() => {
    loadPeriod();
  }, []);

  async function loadPeriod() {
    try {
      setIsLoading(true);
      const data = await adminApi.getContractPeriod();
      setPeriod(data);
      setFormData({
        is_open: data.is_open,
        start_date: data.start_date?.split("T")[0] || "",
        end_date: data.end_date?.split("T")[0] || "",
      });
    } catch (error) {
      toast.error("Failed to load contract period");
    } finally {
      setIsLoading(false);
    }
  }

  async function handleSave() {
    try {
      setIsSaving(true);
      await adminApi.updateContractPeriod({
        is_open: formData.is_open,
        start_date: formData.start_date,
        end_date: formData.end_date,
      });
      toast.success("Contract period updated successfully");
      await loadPeriod();
    } catch (error) {
      toast.error("Failed to update contract period");
    } finally {
      setIsSaving(false);
    }
  }

  if (isLoading) {
    return (
      <>
        <AppHeader />
        <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-2xl">
            <div className="space-y-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <div
                  key={i}
                  className="h-20 animate-pulse rounded-[1.5rem] border border-white/10 bg-white/[0.04]"
                />
              ))}
            </div>
          </div>
        </main>
      </>
    );
  }

  return (
    <>
      <AppHeader />
      <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-2xl space-y-8">
          {/* Header */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
            className="space-y-2"
          >
            <h1 className="text-3xl font-bold text-white">
              Contract Period Management
            </h1>
            <p className="text-zinc-400">
              Control when students can contract courses
            </p>
          </motion.div>

          {/* Status Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <GlassCard className="p-8">
              <div className="space-y-6">
                {/* Toggle */}
                <Toggle
                  label="Contract Period Status"
                  description="Turn ON to allow students to contract courses"
                  checked={formData.is_open}
                  onChange={(checked) =>
                    setFormData({ ...formData, is_open: checked })
                  }
                  disabled={isSaving}
                />

                {/* Current Status Info */}
                <div className="rounded-lg border border-white/10 bg-white/[0.02] p-4">
                  <p className="text-xs uppercase tracking-widest text-zinc-500 mb-3">
                    Current Status
                  </p>
                  <div className="flex items-center gap-3">
                    <div
                      className={[
                        "h-3 w-3 rounded-full",
                        formData.is_open
                          ? "bg-emerald-500"
                          : "bg-red-500/50",
                      ].join(" ")}
                    />
                    <p className="font-medium text-zinc-100">
                      {formData.is_open ? "Contract Period OPEN" : "Contract Period CLOSED"}
                    </p>
                  </div>
                </div>

                <div className="border-t border-white/10" />

                {/* Date Inputs */}
                <div className="space-y-4">
                  <Input
                    type="date"
                    label="Start Date"
                    value={formData.start_date}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        start_date: e.target.value,
                      })
                    }
                    disabled={isSaving}
                  />

                  <Input
                    type="date"
                    label="End Date"
                    value={formData.end_date}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        end_date: e.target.value,
                      })
                    }
                    disabled={isSaving}
                  />
                </div>

                {/* Info */}
                {period && (
                  <div className="rounded-lg border border-blue-500/30 bg-blue-500/[0.1] p-4">
                    <p className="text-xs uppercase tracking-widest text-blue-400 mb-2">
                      Last Updated
                    </p>
                    <p className="text-sm text-blue-200">
                      {new Date(period.updated_at).toLocaleString()}
                    </p>
                  </div>
                )}

                {/* Action Button */}
                <div className="flex gap-3 pt-4">
                  <button
                    onClick={handleSave}
                    disabled={isSaving}
                    className="flex-1 rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-6 py-3 font-medium text-emerald-100 transition hover:border-emerald-500 hover:bg-emerald-500/30 disabled:opacity-50"
                  >
                    {isSaving ? "Saving..." : "Save Changes"}
                  </button>
                </div>
              </div>
            </GlassCard>
          </motion.div>

          {/* Information */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.2 }}
          >
            <GlassCard className="p-6">
              <h3 className="text-sm font-semibold text-zinc-100 mb-3">
                ℹ️ How It Works
              </h3>
              <ul className="space-y-2 text-sm text-zinc-400">
                <li>✓ When ON: Students see the contract system in their dashboard</li>
                <li>✓ When OFF: Students see their schedule and course history</li>
                <li>✓ Set start and end dates to automatically control the period</li>
              </ul>
            </GlassCard>
          </motion.div>
        </div>
      </main>
    </>
  );
}

export default function ContractPeriodPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <div className="flex min-h-screen flex-col">
        <ContractPeriodContent />
      </div>
    </ProtectedRoute>
  );
}
