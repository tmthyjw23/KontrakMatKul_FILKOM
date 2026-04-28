"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { GlassCard } from "@/components/ui/glass-card";
import { AppHeader } from "@/components/layout/app-header";
import { ProtectedRoute } from "@/components/layout/protected-route";
import { DataTable } from "@/components/ui/data-table";
import { Modal, ConfirmDialog } from "@/components/ui/modal";
import { Input } from "@/components/ui/form";
import { adminApi } from "@/lib/api/admin";

interface Student {
  id: string;
  name: string;
  student_number: string;
  email?: string;
}

interface StudentFormData {
  name: string;
  student_number: string;
  email: string;
  password: string;
}

function UsersContent() {
  const [students, setStudents] = useState<Student[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isResetOpen, setIsResetOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [selectedStudent, setSelectedStudent] = useState<Student | null>(null);
  const [resetPassword, setResetPassword] = useState("");
  const [formData, setFormData] = useState<StudentFormData>({
    name: "",
    student_number: "",
    email: "",
    password: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    loadStudents();
  }, []);

  async function loadStudents() {
    try {
      setIsLoading(true);
      const data = await adminApi.getStudents();
      setStudents(data);
    } catch (error) {
      toast.error("Failed to load students");
    } finally {
      setIsLoading(false);
    }
  }

  function openCreateModal() {
    setSelectedStudent(null);
    setFormData({ name: "", student_number: "", email: "", password: "" });
    setErrors({});
    setIsModalOpen(true);
  }

  function validateForm(): boolean {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) newErrors.name = "Name is required";
    if (!formData.student_number.trim())
      newErrors.student_number = "Student number is required";
    if (!formData.password.trim())
      newErrors.password = "Password is required";

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }

  async function handleSave() {
    if (!validateForm()) return;

    try {
      setIsSaving(true);
      await adminApi.createStudent({
        name: formData.name,
        student_number: formData.student_number,
        password: formData.password,
        email: formData.email || undefined,
      });
      toast.success("Student created successfully");
      setIsModalOpen(false);
      await loadStudents();
    } catch (error) {
      toast.error("Failed to create student");
    } finally {
      setIsSaving(false);
    }
  }

  async function handleResetPassword() {
    if (!selectedStudent || !resetPassword.trim()) {
      toast.error("Please enter a new password");
      return;
    }

    try {
      setIsSaving(true);
      await adminApi.resetStudentPassword(selectedStudent.id, resetPassword);
      toast.success("Password reset successfully");
      setIsResetOpen(false);
      setResetPassword("");
    } catch (error) {
      toast.error("Failed to reset password");
    } finally {
      setIsSaving(false);
    }
  }

  async function handleDelete() {
    if (!selectedStudent) return;

    try {
      setIsSaving(true);
      await adminApi.deleteStudent(selectedStudent.id);
      toast.success("Student deleted successfully");
      setIsDeleteOpen(false);
      await loadStudents();
    } catch (error) {
      toast.error("Failed to delete student");
    } finally {
      setIsSaving(false);
    }
  }

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
            className="flex items-center justify-between gap-4"
          >
            <div className="space-y-2">
              <h1 className="text-3xl font-bold text-white">
                User Management
              </h1>
              <p className="text-zinc-400">
                Manage student accounts and credentials
              </p>
            </div>

            <button
              onClick={openCreateModal}
              className="rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-6 py-3 font-medium text-emerald-100 transition hover:border-emerald-500 hover:bg-emerald-500/30"
            >
              + New Student
            </button>
          </motion.div>

          {/* Students Table */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <GlassCard className="p-6">
              <DataTable<Student>
                columns={[
                  { key: "name", label: "Name" },
                  { key: "student_number", label: "Student Number" },
                  {
                    key: "email",
                    label: "Email",
                    render: (email) => email || "—",
                  },
                ]}
                data={students}
                isLoading={isLoading}
                actions={[
                  {
                    label: "Reset Password",
                    onClick: (student) => {
                      setSelectedStudent(student);
                      setIsResetOpen(true);
                    },
                  },
                  {
                    label: "Delete",
                    variant: "danger",
                    onClick: (student) => {
                      setSelectedStudent(student);
                      setIsDeleteOpen(true);
                    },
                  },
                ]}
                emptyMessage="No students found. Create your first student!"
              />
            </GlassCard>
          </motion.div>
        </div>
      </main>

      {/* Create Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => !isSaving && setIsModalOpen(false)}
        title="Create New Student"
        size="md"
        footer={
          <>
            <button
              onClick={() => setIsModalOpen(false)}
              disabled={isSaving}
              className="rounded-lg border border-white/10 bg-white/[0.05] px-4 py-2 text-sm font-medium text-zinc-100 transition hover:bg-white/[0.1] disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              disabled={isSaving}
              className="rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-4 py-2 text-sm font-medium text-emerald-100 transition hover:bg-emerald-500/30 disabled:opacity-50"
            >
              {isSaving ? "Creating..." : "Create Student"}
            </button>
          </>
        }
      >
        <div className="space-y-4">
          <Input
            label="Full Name"
            placeholder="e.g., John Doe"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={errors.name}
            disabled={isSaving}
          />

          <Input
            label="Student Number"
            placeholder="e.g., 22010001"
            value={formData.student_number}
            onChange={(e) =>
              setFormData({ ...formData, student_number: e.target.value })
            }
            error={errors.student_number}
            disabled={isSaving}
          />

          <Input
            label="Email (Optional)"
            type="email"
            placeholder="john@example.com"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            disabled={isSaving}
          />

          <Input
            label="Initial Password"
            type="password"
            placeholder="Set initial password"
            value={formData.password}
            onChange={(e) =>
              setFormData({ ...formData, password: e.target.value })
            }
            error={errors.password}
            disabled={isSaving}
          />
        </div>
      </Modal>

      {/* Reset Password Modal */}
      <Modal
        isOpen={isResetOpen}
        onClose={() => !isSaving && setIsResetOpen(false)}
        title="Reset Password"
        description={`Reset password for ${selectedStudent?.name}`}
        size="sm"
        footer={
          <>
            <button
              onClick={() => setIsResetOpen(false)}
              disabled={isSaving}
              className="rounded-lg border border-white/10 bg-white/[0.05] px-4 py-2 text-sm font-medium text-zinc-100 transition hover:bg-white/[0.1] disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              onClick={handleResetPassword}
              disabled={isSaving}
              className="rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-4 py-2 text-sm font-medium text-emerald-100 transition hover:bg-emerald-500/30 disabled:opacity-50"
            >
              {isSaving ? "Resetting..." : "Reset"}
            </button>
          </>
        }
      >
        <Input
          label="New Password"
          type="password"
          placeholder="Enter new password"
          value={resetPassword}
          onChange={(e) => setResetPassword(e.target.value)}
          disabled={isSaving}
        />
      </Modal>

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={isDeleteOpen}
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteOpen(false)}
        title="Delete Student"
        description={`Are you sure you want to delete student "${selectedStudent?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
        cancelLabel="Cancel"
        isDangerous
        isLoading={isSaving}
      />
    </>
  );
}

export default function UsersPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <div className="flex min-h-screen flex-col">
        <UsersContent />
      </div>
    </ProtectedRoute>
  );
}
