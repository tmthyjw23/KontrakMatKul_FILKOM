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
import type { Course } from "@/src/types/course";

interface CourseFormData {
  code: string;
  name: string;
  sks: string;
  quota: string;
  lecturer: string;
}

function CoursesContent() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [selectedCourse, setSelectedCourse] = useState<Course | null>(null);
  const [formData, setFormData] = useState<CourseFormData>({
    code: "",
    name: "",
    sks: "",
    quota: "",
    lecturer: "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    loadCourses();
  }, []);

  async function loadCourses() {
    try {
      setIsLoading(true);
      const data = await adminApi.getCourses();
      setCourses(data);
    } catch (error) {
      toast.error("Failed to load courses");
    } finally {
      setIsLoading(false);
    }
  }

  function openCreateModal() {
    setSelectedCourse(null);
    setFormData({ code: "", name: "", sks: "", quota: "", lecturer: "" });
    setErrors({});
    setIsModalOpen(true);
  }

  function openEditModal(course: Course) {
    setSelectedCourse(course);
    setFormData({
      code: course.code,
      name: course.name,
      sks: String(course.sks),
      quota: String(course.quota),
      lecturer: course.lecturer,
    });
    setErrors({});
    setIsModalOpen(true);
  }

  function validateForm(): boolean {
    const newErrors: Record<string, string> = {};

    if (!formData.code.trim()) newErrors.code = "Course code is required";
    if (!formData.name.trim()) newErrors.name = "Course name is required";
    if (!formData.sks.trim()) {
      newErrors.sks = "SKS is required";
    } else if (isNaN(Number(formData.sks))) {
      newErrors.sks = "SKS must be a number";
    }
    if (!formData.quota.trim()) {
      newErrors.quota = "Quota is required";
    } else if (isNaN(Number(formData.quota))) {
      newErrors.quota = "Quota must be a number";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }

  async function handleSave() {
    if (!validateForm()) return;

    try {
      setIsSaving(true);
      const payload = {
        code: formData.code,
        name: formData.name,
        sks: Number(formData.sks),
        quota: Number(formData.quota),
        lecturer: formData.lecturer,
      };

      if (selectedCourse) {
        await adminApi.updateCourse(selectedCourse.id, payload);
        toast.success("Course updated successfully");
      } else {
        await adminApi.createCourse(payload);
        toast.success("Course created successfully");
      }

      setIsModalOpen(false);
      await loadCourses();
    } catch (error) {
      toast.error("Failed to save course");
    } finally {
      setIsSaving(false);
    }
  }

  async function handleDelete() {
    if (!selectedCourse) return;

    try {
      setIsSaving(true);
      await adminApi.deleteCourse(selectedCourse.id);
      toast.success("Course deleted successfully");
      setIsDeleteOpen(false);
      await loadCourses();
    } catch (error) {
      toast.error("Failed to delete course");
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
                Course Management
              </h1>
              <p className="text-zinc-400">
                Add, edit, or delete courses from the curriculum
              </p>
            </div>

            <button
              onClick={openCreateModal}
              className="rounded-lg border border-emerald-500/50 bg-emerald-500/20 px-6 py-3 font-medium text-emerald-100 transition hover:border-emerald-500 hover:bg-emerald-500/30"
            >
              + New Course
            </button>
          </motion.div>

          {/* Courses Table */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <GlassCard className="p-6">
              <DataTable<Course>
                columns={[
                  { key: "code", label: "Code", width: "20%" },
                  { key: "name", label: "Name", width: "40%" },
                  {
                    key: "sks",
                    label: "SKS",
                    width: "10%",
                    render: (sks) => <span>{sks}</span>,
                  },
                  {
                    key: "lecturer",
                    label: "Lecturer",
                    width: "30%",
                  },
                ]}
                data={courses}
                isLoading={isLoading}
                actions={[
                  {
                    label: "Edit",
                    onClick: openEditModal,
                  },
                  {
                    label: "Delete",
                    variant: "danger",
                    onClick: (course) => {
                      setSelectedCourse(course);
                      setIsDeleteOpen(true);
                    },
                  },
                ]}
                emptyMessage="No courses found. Create your first course!"
              />
            </GlassCard>
          </motion.div>
        </div>
      </main>

      {/* Create/Edit Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => !isSaving && setIsModalOpen(false)}
        title={selectedCourse ? "Edit Course" : "Create New Course"}
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
              {isSaving ? "Saving..." : "Save Course"}
            </button>
          </>
        }
      >
        <div className="space-y-4">
          <Input
            label="Course Code"
            placeholder="e.g., IFAP272"
            value={formData.code}
            onChange={(e) => setFormData({ ...formData, code: e.target.value })}
            error={errors.code}
            disabled={isSaving}
          />

          <Input
            label="Course Name"
            placeholder="e.g., Automata and Formal Languages"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={errors.name}
            disabled={isSaving}
          />

          <Input
            label="SKS (Credit Hours)"
            type="number"
            placeholder="e.g., 3"
            value={formData.sks}
            onChange={(e) => setFormData({ ...formData, sks: e.target.value })}
            error={errors.sks}
            disabled={isSaving}
          />

          <Input
            label="Quota"
            type="number"
            placeholder="e.g., 40"
            value={formData.quota}
            onChange={(e) => setFormData({ ...formData, quota: e.target.value })}
            error={errors.quota}
            disabled={isSaving}
          />

          <Input
            label="Lecturer Name"
            placeholder="e.g., Dr. John Doe"
            value={formData.lecturer}
            onChange={(e) =>
              setFormData({ ...formData, lecturer: e.target.value })
            }
            disabled={isSaving}
          />
        </div>
      </Modal>

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={isDeleteOpen}
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteOpen(false)}
        title="Delete Course"
        description={`Are you sure you want to delete "${selectedCourse?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
        cancelLabel="Cancel"
        isDangerous
        isLoading={isSaving}
      />
    </>
  );
}

export default function CoursesPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <div className="flex min-h-screen flex-col">
        <CoursesContent />
      </div>
    </ProtectedRoute>
  );
}
