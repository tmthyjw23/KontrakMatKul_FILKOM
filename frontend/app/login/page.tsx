"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { motion } from "framer-motion";
import { Input } from "@/components/ui/form";
import { useAuthStore } from "@/lib/store/useAuthStore";
import { authApi } from "@/lib/api/admin";
import type { UserRole } from "@/src/types/auth";

export default function LoginPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  const [role, setRole] = useState<UserRole>("student");
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    student_number: "",
    password: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  function validateForm(): boolean {
    const newErrors: Record<string, string> = {};

    if (role === "student" && !formData.student_number.trim()) {
      newErrors.student_number = "Student number is required";
    }
    if (!formData.password.trim()) {
      newErrors.password = "Password is required";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }

  async function handleLogin(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      const response = await authApi.login({
        student_number: role === "student" ? formData.student_number : "admin",
        password: formData.password,
        role,
      });

      setAuth(response.user, response.token);
      toast.success(`Welcome, ${response.user.name}!`);

      // Redirect based on role
      router.push(role === "admin" ? "/admin" : "/student");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Login failed. Please try again.";
      toast.error(message);
      setErrors({ submit: message });
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <main className="min-h-screen flex items-center justify-center px-4 py-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md"
      >
        <div className="rounded-[1.75rem] border border-white/10 bg-white/[0.05] p-8 backdrop-blur-2xl">
          {/* Logo & Title */}
          <div className="mb-8 text-center">
            <h1 className="text-3xl font-bold text-white">FILKOM</h1>
            <p className="mt-2 text-sm text-zinc-400">
              Course Registration System
            </p>
          </div>

          {/* Role Selection */}
          <div className="mb-8 space-y-3">
            <p className="text-sm font-medium text-zinc-200">Login As</p>
            <div className="grid grid-cols-2 gap-3">
              {[
                { value: "student", label: "Student" },
                { value: "admin", label: "Admin" },
              ].map((option) => (
                <button
                  key={option.value}
                  onClick={() => {
                    setRole(option.value as UserRole);
                    setErrors({});
                  }}
                  className={[
                    "rounded-lg border px-4 py-3 text-sm font-medium transition",
                    role === option.value
                      ? "border-emerald-500/50 bg-emerald-500/20 text-emerald-100"
                      : "border-white/10 bg-white/[0.05] text-zinc-300 hover:border-white/20",
                  ].join(" ")}
                >
                  {option.label}
                </button>
              ))}
            </div>
          </div>

          {/* Form */}
          <form onSubmit={handleLogin} className="space-y-4">
            {role === "student" && (
              <Input
                type="text"
                label="Student Number"
                placeholder="e.g., 22010001"
                value={formData.student_number}
                onChange={(e) =>
                  setFormData({ ...formData, student_number: e.target.value })
                }
                error={errors.student_number}
                disabled={isLoading}
              />
            )}

            <Input
              type="password"
              label="Password"
              placeholder="Enter your password"
              value={formData.password}
              onChange={(e) =>
                setFormData({ ...formData, password: e.target.value })
              }
              error={errors.password}
              disabled={isLoading}
            />

            {errors.submit && (
              <div className="rounded-lg border border-red-500/30 bg-red-500/[0.1] p-3">
                <p className="text-xs text-red-300">{errors.submit}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="w-full rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-4 py-3 font-medium text-emerald-100 transition hover:border-emerald-500 hover:bg-emerald-500/30 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? "Logging in..." : "Login"}
            </button>
          </form>

          {/* Demo Credentials */}
          <div className="mt-8 rounded-lg border border-white/10 bg-white/[0.02] p-4">
            <p className="text-xs font-medium text-zinc-400 uppercase tracking-widest">
              Demo Credentials
            </p>
            <div className="mt-3 space-y-2 text-xs text-zinc-400">
              <div>
                <p className="font-medium text-zinc-300">Student</p>
                <p>Number: 22010001</p>
                <p>Password: password123</p>
              </div>
              <div className="border-t border-white/10 pt-2">
                <p className="font-medium text-zinc-300">Admin</p>
                <p>Password: admin123</p>
              </div>
            </div>
          </div>
        </div>
      </motion.div>
    </main>
  );
}
