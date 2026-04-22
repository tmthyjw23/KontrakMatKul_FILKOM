"use client";

import { motion, AnimatePresence } from "framer-motion";
import type { ReactNode } from "react";

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children?: ReactNode;
  footer?: ReactNode;
  size?: "sm" | "md" | "lg";
}

const sizeClasses = {
  sm: "max-w-sm",
  md: "max-w-md",
  lg: "max-w-lg",
};

export function Modal({
  isOpen,
  onClose,
  title,
  description,
  children,
  footer,
  size = "md",
}: ModalProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
          />
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            className={`fixed left-1/2 top-1/2 z-50 w-full -translate-x-1/2 -translate-y-1/2 ${sizeClasses[size]} rounded-[1.75rem] border border-white/10 bg-white/[0.05] p-6 sm:p-8 shadow-2xl backdrop-blur-2xl`}
          >
            <div className="space-y-4">
              <div>
                <h2 className="text-2xl font-semibold text-white">{title}</h2>
                {description && (
                  <p className="mt-2 text-sm text-zinc-400">{description}</p>
                )}
              </div>

              {children && <div className="space-y-4">{children}</div>}

              {footer && (
                <div className="flex gap-3 justify-end border-t border-white/10 pt-4">
                  {footer}
                </div>
              )}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}

interface ConfirmDialogProps {
  isOpen: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  title: string;
  description?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  isDangerous?: boolean;
  isLoading?: boolean;
}

export function ConfirmDialog({
  isOpen,
  onConfirm,
  onCancel,
  title,
  description,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  isDangerous = false,
  isLoading = false,
}: ConfirmDialogProps) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onCancel}
      title={title}
      description={description}
      size="sm"
      footer={
        <>
          <button
            onClick={onCancel}
            disabled={isLoading}
            className="rounded-lg border border-white/10 bg-white/[0.05] px-4 py-2 text-sm font-medium text-zinc-100 transition hover:bg-white/[0.1] disabled:opacity-50"
          >
            {cancelLabel}
          </button>
          <button
            onClick={onConfirm}
            disabled={isLoading}
            className={[
              "rounded-lg px-4 py-2 text-sm font-medium transition disabled:opacity-50",
              isDangerous
                ? "border border-red-500/50 bg-red-500/20 text-red-200 hover:bg-red-500/30"
                : "border border-emerald-500/50 bg-emerald-500/20 text-emerald-200 hover:bg-emerald-500/30",
            ].join(" ")}
          >
            {isLoading ? "Loading..." : confirmLabel}
          </button>
        </>
      }
    />
  );
}
