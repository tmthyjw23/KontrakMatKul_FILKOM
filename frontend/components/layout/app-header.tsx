"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth, useAuthStore, useUserRole } from "@/lib/store/useAuthStore";

export function AppHeader() {
  const router = useRouter();
  const { user } = useAuth();
  const userRole = useUserRole();
  const clearAuth = useAuthStore((state) => state.clearAuth);

  function handleLogout() {
    clearAuth();
    router.push("/login");
  }

  if (!user) {
    return null;
  }

  return (
    <header className="sticky top-0 z-30 border-b border-white/10 bg-black/40 backdrop-blur-md">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-8">
          <Link href={userRole === "admin" ? "/admin" : "/student"} className="flex items-center gap-2">
            <div className="text-xl font-bold text-white">FILKOM</div>
          </Link>

          <nav className="hidden sm:flex gap-6">
            {userRole === "admin" ? (
              <>
                <Link
                  href="/admin"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Dashboard
                </Link>
                <Link
                  href="/admin/contract-period"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Periode
                </Link>
                <Link
                  href="/admin/courses"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Kurikulum
                </Link>
                <Link
                  href="/admin/monitoring"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Monitoring
                </Link>
                <Link
                  href="/admin/users"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Users
                </Link>
              </>
            ) : (
              <>
                <Link
                  href="/student"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Dashboard
                </Link>
                <Link
                  href="/student/schedule"
                  className="text-sm text-zinc-400 hover:text-zinc-100 transition"
                >
                  Jadwal
                </Link>
              </>
            )}
          </nav>
        </div>

        <div className="flex items-center gap-4">
          <div className="hidden sm:block text-right">
            <p className="text-sm font-medium text-zinc-100">{user.name}</p>
            <p className="text-xs text-zinc-500">
              {userRole === "admin" ? "Admin" : user.student_number}
            </p>
          </div>

          <button
            onClick={handleLogout}
            className="rounded-lg border border-red-500/30 bg-red-500/[0.12] px-4 py-2 text-sm font-medium text-red-100 transition hover:border-red-500/50 hover:bg-red-500/20"
          >
            Logout
          </button>
        </div>
      </div>
    </header>
  );
}
