"use client";
import { useState } from "react";
import useSWR from "swr";
import { useRouter, useSearchParams } from "next/navigation";
import Money from "@/components/Money";
import StatusBadge from "@/components/StatusBadge";
import Toast from "@/components/Toast";
import { api } from "@/lib/api";

interface Receivable {
  id: string;
  dueDate: string;
  amount: number;
  status: string;
  customer: { trade_name: string };
  service: { name: string };
}

const fetcher = (url: string) => api<Receivable[]>(url);

export default function ReceivablesClient() {
  const router = useRouter();
  const search = useSearchParams();
  const statusParam = search.get("status") || "open";
  const { data, isLoading, mutate } = useSWR(
    `/receivables?status=${statusParam}`,
    fetcher,
  );
  const [toast, setToast] = useState<string | null>(null);

  const changeStatus = (s: string) => {
    const params = new URLSearchParams(search.toString());
    if (s) params.set("status", s);
    else params.delete("status");
    router.replace(`/finance/receivables?${params.toString()}`);
  };

  const markPaid = async (id: string) => {
    try {
      await api(`/receivables/${id}/pay`, { method: "PUT" });
      setToast("Marcado como pago");
      mutate();
    } catch {
      setToast("Erro ao marcar como pago");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  const tabs = [
    { label: "Todos", value: "" },
    { label: "Abertos", value: "open" },
    { label: "Pagos", value: "paid" },
    { label: "Vencidos", value: "overdue" },
  ];

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Contas a Receber</h1>
      <div className="flex gap-2 mb-4">
        {tabs.map((t) => (
          <button
            key={t.label}
            onClick={() => changeStatus(t.value)}
            className={`px-3 py-1 border rounded ${
              statusParam === (t.value || "open") ? "bg-primary text-background" : ""
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>
      {isLoading ? (
        <div className="flex justify-center p-4">
          <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-card">
            <thead className="bg-card">
              <tr>
                <th className="px-4 py-2 text-left">Vencimento</th>
                <th className="px-4 py-2 text-left">Cliente</th>
                <th className="px-4 py-2 text-left">Serviço</th>
                <th className="px-4 py-2 text-left">Valor</th>
                <th className="px-4 py-2 text-left">Status</th>
                <th className="px-4 py-2 text-left">Ações</th>
              </tr>
            </thead>
            <tbody>
              {data?.map((r, idx) => (
                <tr key={r.id} className={idx % 2 ? "bg-card" : "bg-card"}>
                  <td className="px-4 py-2">{new Date(r.dueDate).toLocaleDateString()}</td>
                  <td className="px-4 py-2">{r.customer?.trade_name}</td>
                  <td className="px-4 py-2">{r.service?.name}</td>
                  <td className="px-4 py-2"><Money value={r.amount} /></td>
                  <td className="px-4 py-2"><StatusBadge status={r.status} /></td>
                  <td className="px-4 py-2">
                    {r.status === "open" && (
                      <button className="text-primary" onClick={() => markPaid(r.id)}>
                        Marcar como Pago
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <Toast message={toast} />
    </div>
  );
}
