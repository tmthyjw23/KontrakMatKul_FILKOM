"use client";

import { GlassCard } from "./glass-card";

interface DataTableColumn<T> {
  key: keyof T;
  label: string;
  render?: (value: T[keyof T], item: T) => React.ReactNode;
  width?: string;
}

interface DataTableProps<T extends { id?: string | number }> {
  columns: DataTableColumn<T>[];
  data: T[];
  isLoading?: boolean;
  onRowClick?: (item: T) => void;
  actions?: Array<{
    label: string;
    onClick: (item: T) => void;
    variant?: "default" | "danger" | "success";
  }>;
  emptyMessage?: string;
}

export function DataTable<T extends { id?: string | number }>({
  columns,
  data,
  isLoading = false,
  onRowClick,
  actions,
  emptyMessage = "No data available",
}: DataTableProps<T>) {
  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <div
            key={i}
            className="h-12 animate-pulse rounded-lg border border-white/10 bg-white/[0.04]"
          />
        ))}
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <GlassCard className="p-8 text-center">
        <p className="text-sm text-zinc-400">{emptyMessage}</p>
      </GlassCard>
    );
  }

  return (
    <div className="overflow-x-auto rounded-[1.75rem] border border-white/10">
      <table className="w-full text-sm">
        <thead className="border-b border-white/10 bg-white/[0.03]">
          <tr>
            {columns.map((col) => (
              <th
                key={String(col.key)}
                className="px-6 py-4 text-left font-medium text-zinc-400"
                style={{ width: col.width }}
              >
                {col.label}
              </th>
            ))}
            {actions && actions.length > 0 && (
              <th className="px-6 py-4 text-center font-medium text-zinc-400">
                Actions
              </th>
            )}
          </tr>
        </thead>
        <tbody className="divide-y divide-white/10">
          {data.map((item, idx) => (
            <tr
              key={item.id ?? idx}
              onClick={() => onRowClick?.(item)}
              className={
                onRowClick ? "cursor-pointer hover:bg-white/[0.02] transition" : ""
              }
            >
              {columns.map((col) => (
                <td key={String(col.key)} className="px-6 py-4 text-zinc-100">
                  {col.render
                    ? col.render(item[col.key], item)
                    : String(item[col.key] ?? "")}
                </td>
              ))}
              {actions && actions.length > 0 && (
                <td className="px-6 py-4 text-center">
                  <div className="flex gap-2 justify-center">
                    {actions.map((action) => (
                      <button
                        key={action.label}
                        onClick={(e) => {
                          e.stopPropagation();
                          action.onClick(item);
                        }}
                        className={[
                          "px-3 py-1.5 text-xs font-medium rounded-lg transition",
                          action.variant === "danger"
                            ? "bg-red-500/20 text-red-200 hover:bg-red-500/30"
                            : action.variant === "success"
                              ? "bg-emerald-500/20 text-emerald-200 hover:bg-emerald-500/30"
                              : "bg-white/[0.1] text-zinc-300 hover:bg-white/[0.15]",
                        ].join(" ")}
                      >
                        {action.label}
                      </button>
                    ))}
                  </div>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
