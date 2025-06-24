"use client";
import ProtectedRoute from "@/components/ProtectedRoute";
import Link from "next/link";
import { getToken } from "@/hooks/useAuth";

function getRole(): string | null {
  const token = getToken();
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.role || null;
  } catch {
    return null;
  }
}

export default function DashboardPage() {
  const role = getRole();
  return (
    <ProtectedRoute>
      <div className="p-4 space-y-2">
        <div>Dashboard</div>
        {role === "admin" && (
          <Link href="/audit-logs" className="text-blue-600 hover:underline">
            Audit Logs
          </Link>
        )}
      </div>
    </ProtectedRoute>
  );
}
