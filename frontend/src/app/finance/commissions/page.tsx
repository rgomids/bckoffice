"use client";
import { useState } from "react";
import useSWR from "swr";
import ProtectedRoute from "@/components/ProtectedRoute";
import Money from "@/components/Money";
import Toast from "@/components/Toast";
import { api } from "@/util/api";

interface Commission {
  id: string;
  contract: { id: string };
  promoter: { name: string };
  amount: number;
  approved: boolean;
}

const fetcher = (url: string) => api<Commission[]>(url);

export default function CommissionsPage() {
  const { data, isLoading, mutate } = useSWR(
    "/commissions?pending=true",
    fetcher,
  );
  const [toast, setToast] = useState<string | null>(null);

  const approve = async (id: string) => {
    try {
      await api(`/commissions/${id}/approve`, { method: "PUT" });
      setToast("Comissão aprovada");
      mutate();
    } catch {
      setToast("Erro ao aprovar comissão");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  return (
    <ProtectedRoute roles={["finance", "admin"]}>
      <div className="p-4">
        <h1 className="text-xl font-bold mb-4">Comissões</h1>
        {isLoading ? (
          <div className="flex justify-center p-4">
            <div className="h-5 w-5 border-2 border-gray-300 border-t-transparent rounded-full animate-spin" />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left">Contrato</th>
                  <th className="px-4 py-2 text-left">Promotor</th>
                  <th className="px-4 py-2 text-left">Valor</th>
                  <th className="px-4 py-2 text-left">Approved</th>
                  <th className="px-4 py-2 text-left">Ações</th>
                </tr>
              </thead>
              <tbody>
                {data?.map((c, idx) => (
                  <tr key={c.id} className={idx % 2 ? "bg-gray-50" : "bg-white"}>
                    <td className="px-4 py-2">{c.contract?.id}</td>
                    <td className="px-4 py-2">{c.promoter?.name}</td>
                    <td className="px-4 py-2"><Money value={c.amount} /></td>
                    <td className="px-4 py-2">{c.approved ? "Yes" : "No"}</td>
                    <td className="px-4 py-2">
                      {!c.approved && (
                        <button className="text-blue-600" onClick={() => approve(c.id)}>
                          Aprovar
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
    </ProtectedRoute>
  );
}
