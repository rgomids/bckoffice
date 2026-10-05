"use client";
import useSWR from "swr";
import Link from "next/link";
import ProtectedRoute from "@/components/ProtectedRoute";
import { api } from "@/lib/api";

interface AuditLog {
  id: string;
  entityName: string;
  entityId: string;
  action: string;
  createdAt: string;
  userId: string;
}

const fetcher = (url: string) => api<AuditLog[]>(url);

export default function AuditLogsPage() {
  const { data } = useSWR("/audit-logs?limit=50", fetcher);
  return (
    <ProtectedRoute roles={["admin"]}>
      <div className="p-4 space-y-2">
        <Link href="/dashboard" className="text-primary hover:underline">
          Voltar
        </Link>
        <h1 className="text-xl font-bold">Audit Logs</h1>
        {data ? (
          <div className="space-y-1 font-mono text-sm">
            {data.map((l) => (
              <div key={l.id}>{JSON.stringify(l)}</div>
            ))}
          </div>
        ) : (
          <div>Carregando...</div>
        )}
      </div>
    </ProtectedRoute>
  );
}
