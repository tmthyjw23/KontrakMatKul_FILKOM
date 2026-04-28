"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { DataTable } from "@/components/ui/data-table";
import { adminApi } from "@/lib/api/admin";
import type { StudentEnrollment } from "@/src/types/auth";

function MonitoringContent() {
  const [enrollments, setEnrollments] = useState<StudentEnrollment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isProcessing, setIsProcessing] = useState(false);
  const [filter, setFilter] = useState<"all" | "pending" | "approved" | "rejected">("all");

  useEffect(() => {
    loadEnrollments();
  }, []);

  async function loadEnrollments() {
    try {
      setIsLoading(true);
      const data = await adminApi.getEnrollments();
      setEnrollments(data);
    } catch (error) {
      toast.error("Failed to load enrollments");
    } finally {
      setIsLoading(false);
    }
  }

  async function handleApprove(enrollment: StudentEnrollment) {
    try {
      setIsProcessing(true);
      await adminApi.approveEnrollment(enrollment.id);
      toast.success(`Approved enrollment for ${enrollment.student_name}`);
      await loadEnrollments();
    } catch (error) {
      toast.error("Failed to approve enrollment");
    } finally {
      setIsProcessing(false);
    }
  }

  async function handleReject(enrollment: StudentEnrollment) {
    try {
      setIsProcessing(true);
      await adminApi.rejectEnrollment(enrollment.id);
      toast.success(`Rejected enrollment for ${enrollment.student_name}`);
      await loadEnrollments();
    } catch (error) {
      toast.error("Failed to reject enrollment");
    } finally {
      setIsProcessing(false);
    }
  }

  const filteredEnrollments =
    filter === "all"
      ? enrollments
      : enrollments.filter((e) => e.status === filter);

  const stats = {
    total: enrollments.length,
    pending: enrollments.filter((e) => e.status === "pending").length,
    approved: enrollments.filter((e) => e.status === "approved").length,
    rejected: enrollments.filter((e) => e.status === "rejected").length,
  };

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
            <h1 className="text-3xl font-bold text-white">
              Contract Monitoring
            </h1>
            <p className="text-zinc-400">
              View and manage student course enrollments
            </p>
          </motion.div>

          {/* Statistics */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
            className="grid gap-4 sm:grid-cols-4"
          >
            {[
              { label: "Total Enrollments", value: stats.total, color: "blue" },
              { label: "Pending", value: stats.pending, color: "orange" },
              { label: "Approved", value: stats.approved, color: "emerald" },
              { label: "Rejected", value: stats.rejected, color: "red" },
            ].map((stat) => (
              <GlassCard key={stat.label} className="p-6">
                <p className="text-xs uppercase tracking-widest text-zinc-500 mb-2">
                  {stat.label}
                </p>
                <p className="text-3xl font-bold text-white">{stat.value}</p>
              </GlassCard>
            ))}
          </motion.div>

          {/* Filters */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.2 }}
            className="flex gap-2 overflow-x-auto pb-2"
          >
            {[
              { value: "all", label: "All Enrollments" },
              { value: "pending", label: "Pending" },
              { value: "approved", label: "Approved" },
              { value: "rejected", label: "Rejected" },
            ].map((f) => (
              <button
                key={f.value}
                onClick={() => setFilter(f.value as typeof filter)}
                className={[
                  "px-4 py-2 rounded-lg text-sm font-medium transition whitespace-nowrap",
                  filter === f.value
                    ? "border-emerald-500/50 bg-emerald-500/20 text-emerald-100"
                    : "border border-white/10 text-zinc-300 hover:border-white/20",
                ].join(" ")}
              >
                {f.label}
              </button>
            ))}
          </motion.div>

          {/* Enrollments Table */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.3 }}
          >
            <GlassCard className="p-6">
              <DataTable<StudentEnrollment>
                columns={[
                  { key: "student_name", label: "Student Name" },
                  { key: "student_number", label: "Student Number" },
                  { key: "course_code", label: "Course Code" },
                  { key: "course_name", label: "Course Name" },
                  {
                    key: "status",
                    label: "Status",
                    render: (status) => (
                      <span
                        className={[
                          "px-3 py-1 rounded-full text-xs font-medium",
                          status === "pending"
                            ? "bg-orange-500/20 text-orange-200"
                            : status === "approved"
                              ? "bg-emerald-500/20 text-emerald-200"
                              : "bg-red-500/20 text-red-200",
                        ].join(" ")}
                      >
                        {String(status).charAt(0).toUpperCase() +
                          String(status).slice(1)}
                      </span>
                    ),
                  },
                ]}
                data={filteredEnrollments}
                isLoading={isLoading}
                actions={
                  filter === "pending"
                    ? [
                        {
                          label: "Approve",
                          variant: "success",
                          onClick: handleApprove,
                        },
                        {
                          label: "Reject",
                          variant: "danger",
                          onClick: handleReject,
                        },
                      ]
                    : undefined
                }
                emptyMessage="No enrollments found"
              />
            </GlassCard>
          </motion.div>
        </div>
      </main>
    </>
  );
}

export default function MonitoringPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <div className="flex min-h-screen flex-col">
        <MonitoringContent />
      </div>
    </ProtectedRoute>
  );
}
