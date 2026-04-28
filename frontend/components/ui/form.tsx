"use client";

import type { InputHTMLAttributes, SelectHTMLAttributes, TextareaHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
}

export function Input({
  label,
  error,
  helperText,
  className,
  ...props
}: InputProps) {
  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-zinc-200">
          {label}
        </label>
      )}
      <input
        className={[
          "w-full rounded-lg border bg-white/[0.03] px-4 py-2.5 text-sm text-zinc-100 placeholder-zinc-500 transition",
          error
            ? "border-red-500/50 focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500/30"
            : "border-white/10 focus:border-white/30 focus:outline-none focus:ring-1 focus:ring-white/20",
          className,
        ].join(" ")}
        {...props}
      />
      {error && <p className="text-xs text-red-400">{error}</p>}
      {helperText && !error && <p className="text-xs text-zinc-500">{helperText}</p>}
    </div>
  );
}

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  options: Array<{ value: string; label: string }>;
}

export function Select({
  label,
  error,
  options,
  className,
  ...props
}: SelectProps) {
  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-zinc-200">
          {label}
        </label>
      )}
      <select
        className={[
          "w-full rounded-lg border bg-white/[0.03] px-4 py-2.5 text-sm text-zinc-100 transition",
          error
            ? "border-red-500/50 focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500/30"
            : "border-white/10 focus:border-white/30 focus:outline-none focus:ring-1 focus:ring-white/20",
          className,
        ].join(" ")}
        {...props}
      >
        <option value="">Select an option</option>
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && <p className="text-xs text-red-400">{error}</p>}
    </div>
  );
}

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
}

export function Textarea({
  label,
  error,
  className,
  ...props
}: TextareaProps) {
  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-zinc-200">
          {label}
        </label>
      )}
      <textarea
        className={[
          "w-full rounded-lg border bg-white/[0.03] px-4 py-2.5 text-sm text-zinc-100 placeholder-zinc-500 transition",
          error
            ? "border-red-500/50 focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500/30"
            : "border-white/10 focus:border-white/30 focus:outline-none focus:ring-1 focus:ring-white/20",
          className,
        ].join(" ")}
        {...props}
      />
      {error && <p className="text-xs text-red-400">{error}</p>}
    </div>
  );
}

interface ToggleProps {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  description?: string;
  disabled?: boolean;
}

export function Toggle({
  label,
  checked,
  onChange,
  description,
  disabled = false,
}: ToggleProps) {
  return (
    <div className="flex items-center justify-between py-4">
      <div>
        <p className="font-medium text-zinc-100">{label}</p>
        {description && <p className="text-sm text-zinc-400">{description}</p>}
      </div>
      <button
        onClick={() => !disabled && onChange(!checked)}
        disabled={disabled}
        className={[
          "relative inline-flex h-8 w-14 items-center rounded-full transition",
          checked ? "bg-emerald-500/30 border border-emerald-500/50" : "bg-white/[0.1] border border-white/10",
          disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer",
        ].join(" ")}
      >
        <span
          className={[
            "inline-block h-6 w-6 transform rounded-full bg-white transition",
            checked ? "translate-x-7" : "translate-x-1",
          ].join(" ")}
        />
      </button>
    </div>
  );
}
