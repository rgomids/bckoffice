"use client";

interface StatusBadgeProps {
  status: string;
}

const COLORS: Record<string, string> = {
  paid: "bg-green-200 text-green-800",
  open: "bg-yellow-200 text-yellow-800",
  overdue: "bg-red-200 text-red-800",
};

export default function StatusBadge({ status }: StatusBadgeProps) {
  const cls = COLORS[status] || "bg-card text-foreground";
  return <span className={`px-2 py-1 rounded text-xs ${cls}`}>{status}</span>;
}
