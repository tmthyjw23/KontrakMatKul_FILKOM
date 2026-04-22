"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth, useAuthStore } from "@/lib/store/useAuthStore";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";

// Wrap the original dashboard
import Home from "@/app/page";

function StudentDashboardContent() {
  const { user } = useAuth();
  const router = useRouter();
  const loadFromStorage = useAuthStore((state) => state.loadFromStorage);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  // TODO: Fetch contract period status and conditionally render
  // For now, always show the contract system
  const contractPeriodIsOpen = true;

  if (!user) {
    return null;
  }

  return (
    <>
      <AppHeader />
      <Home />
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
