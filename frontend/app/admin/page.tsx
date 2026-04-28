"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";

const adminMenus = [
  {
    title: "Contract Period Management",
    description: "Configure contract period status, start and end dates",
    href: "/admin/contract-period",
    icon: "📅",
    color: "from-blue-500/20 to-cyan-500/10",
  },
  {
    title: "Course Management",
    description: "Add, edit, delete courses and set prerequisites",
    href: "/admin/courses",
    icon: "📚",
    color: "from-purple-500/20 to-pink-500/10",
  },
  {
    title: "Contract Monitoring",
    description: "View and approve student course enrollments",
    href: "/admin/monitoring",
    icon: "👥",
    color: "from-emerald-500/20 to-teal-500/10",
  },
  {
    title: "User Management",
    description: "Manage student accounts, passwords, and data",
    href: "/admin/users",
    icon: "👤",
    color: "from-orange-500/20 to-red-500/10",
  },
];

function AdminDashboardContent() {
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
            <h1 className="text-4xl font-bold text-white">Admin Dashboard</h1>
            <p className="text-zinc-400">
              Manage contracts, courses, and student data
            </p>
          </motion.div>

          {/* Menu Grid */}
          <div className="grid gap-4 sm:grid-cols-2">
            {adminMenus.map((menu, idx) => (
              <motion.div
                key={menu.href}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4, delay: idx * 0.1 }}
              >
                <Link href={menu.href}>
                  <GlassCard className={`h-full p-6 sm:p-8 cursor-pointer transition hover:border-white/20 bg-gradient-to-br ${menu.color}`}>
                    <div className="space-y-4">
                      <div className="text-4xl">{menu.icon}</div>
                      <div>
                        <h3 className="text-lg font-semibold text-white">
                          {menu.title}
                        </h3>
                        <p className="mt-2 text-sm text-zinc-400">
                          {menu.description}
                        </p>
                      </div>
                      <div className="flex items-center text-sm text-zinc-400 group">
                        Go to page <span className="ml-2 group-hover:translate-x-1 transition">→</span>
                      </div>
                    </div>
                  </GlassCard>
                </Link>
              </motion.div>
            ))}
          </div>

          {/* Quick Stats (Placeholder for future) */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.4 }}
          >
            <GlassCard className="p-8">
              <h2 className="text-xl font-semibold text-white mb-6">
                Quick Statistics
              </h2>
              <div className="grid gap-6 sm:grid-cols-4">
                {[
                  { label: "Active Students", value: "—" },
                  { label: "Total Courses", value: "—" },
                  { label: "Pending Enrollments", value: "—" },
                  { label: "Contract Status", value: "—" },
                ].map((stat) => (
                  <div key={stat.label} className="space-y-2">
                    <p className="text-xs uppercase tracking-widest text-zinc-500">
                      {stat.label}
                    </p>
                    <p className="text-2xl font-bold text-white">
                      {stat.value}
                    </p>
                  </div>
                ))}
              </div>
            </GlassCard>
          </motion.div>
        </div>
      </main>
    </>
  );
}

export default function AdminDashboard() {
  return (
    <ProtectedRoute requiredRole="admin">
      <div className="flex min-h-screen flex-col">
        <AdminDashboardContent />
      </div>
    </ProtectedRoute>
  );
}
