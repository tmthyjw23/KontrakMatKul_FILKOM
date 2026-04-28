"use client";

import { useEffect, useState } from "react";
import { useAuth, useAuthStore } from "@/lib/store/useAuthStore";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { GlassCard } from "@/components/ui/glass-card";
import { studentApi } from "@/lib/api/admin";
import type { ContractPeriod } from "@/src/types/auth";

import Home from "@/app/page";

function StudentDashboardContent() {
  const { user } = useAuth();
  const loadFromStorage = useAuthStore((state) => state.loadFromStorage);
  const [contractPeriod, setContractPeriod] = useState<ContractPeriod | null>(
    null
  );
  const [periodLoading, setPeriodLoading] = useState(true);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  useEffect(() => {
    studentApi
      .getContractPeriod()
      .then(setContractPeriod)
      .catch(() => {
        // Endpoint not yet available — default to open so students aren't locked out
        setContractPeriod({
          id: "",
          is_open: true,
          start_date: "",
          end_date: "",
          created_at: "",
          updated_at: "",
        });
      })
      .finally(() => setPeriodLoading(false));
  }, []);

  if (!user) return null;

  // While fetching, default to open to avoid a flash of the "closed" screen
  const contractPeriodIsOpen = periodLoading
    ? true
    : (contractPeriod?.is_open ?? true);

  return (
    <>
      <AppHeader />
      {contractPeriodIsOpen ? (
        <Home />
      ) : (
        <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-2xl py-16">
            <GlassCard className="p-12 text-center">
              <p className="text-4xl mb-4">🔒</p>
              <h2 className="text-2xl font-bold text-white mb-3">
                Contract Period Closed
              </h2>
              <p className="text-zinc-400 leading-relaxed">
                Course registration is currently closed. Please check back when
                the contract period opens.
              </p>
              {contractPeriod?.start_date && (
                <p className="mt-4 text-sm text-zinc-500">
                  Opens:{" "}
                  {new Date(contractPeriod.start_date).toLocaleDateString(
                    "id-ID",
                    { day: "numeric", month: "long", year: "numeric" }
                  )}
                </p>
              )}
            </GlassCard>
          </div>
        </main>
      )}
    </>
  );
}

export default function StudentDashboardPage() {
  return (
    <ProtectedRoute requiredRole="student">
      <div className="flex min-h-screen flex-col">
        <StudentDashboardContent />
      </div>
    </ProtectedRoute>
  );
}
