"use client";

import { useAuth } from "@/lib/store/useAuthStore";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { GlassCard } from "@/components/ui/glass-card";
import { motion } from "framer-motion";

function StudentProfileContent() {
  const { user } = useAuth();

  if (!user) {
    return null;
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
            <h1 className="text-3xl font-bold text-white">My Profile</h1>
            <p className="text-zinc-400">View and manage your account information</p>
          </motion.div>

          {/* Profile Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <GlassCard className="p-8">
              <div className="space-y-6">
                {/* Avatar */}
                <div className="flex items-center gap-4">
                  <div className="h-16 w-16 rounded-full bg-gradient-to-br from-emerald-400 to-cyan-400 flex items-center justify-center text-2xl font-bold text-black">
                    {user.name.charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <p className="text-xs uppercase tracking-widest text-zinc-500">
                      Student
                    </p>
                    <h2 className="text-2xl font-bold text-white mt-1">
                      {user.name}
                    </h2>
                  </div>
                </div>

                <div className="border-t border-white/10" />

                {/* Information */}
                <div className="space-y-4">
                  <div>
                    <p className="text-xs uppercase tracking-widest text-zinc-500 mb-2">
                      Student Number
                    </p>
                    <p className="text-lg font-semibold text-zinc-100">
                      {user.student_number || "—"}
                    </p>
                  </div>

                  <div>
                    <p className="text-xs uppercase tracking-widest text-zinc-500 mb-2">
                      Email
                    </p>
                    <p className="text-lg font-semibold text-zinc-100">
                      {user.email || "Not set"}
                    </p>
                  </div>

                  <div>
                    <p className="text-xs uppercase tracking-widest text-zinc-500 mb-2">
                      Role
                    </p>
                    <p className="text-lg font-semibold text-emerald-100">
                      Student
                    </p>
                  </div>
                </div>
              </div>
            </GlassCard>
          </motion.div>

          {/* Help Section */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.2 }}
          >
            <GlassCard className="p-6">
              <h3 className="text-sm font-semibold text-zinc-100 mb-3">
                📋 Quick Links
              </h3>
              <ul className="space-y-2 text-sm text-zinc-400">
                <li>
                  • Go to{" "}
                  <span className="text-emerald-200">Dashboard</span> to view
                  courses
                </li>
                <li>
                  • Go to{" "}
                  <span className="text-emerald-200">Jadwal</span> to see your
                  schedule
                </li>
                <li>
                  • Contact your department if you need to change password
                </li>
              </ul>
            </GlassCard>
          </motion.div>
        </div>
      </main>
    </>
  );
}

export default function StudentProfilePage() {
  return (
    <ProtectedRoute requiredRole="student">
      <div className="flex min-h-screen flex-col">
        <StudentProfileContent />
      </div>
    </ProtectedRoute>
  );
}
