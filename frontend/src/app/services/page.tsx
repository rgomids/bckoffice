"use client";

import { useState } from "react";
import { Dialog, Transition } from "@headlessui/react";
import useSWR from "swr";
import ProtectedRoute from "@/components/ProtectedRoute";
import Money from "@/components/Money";
import Toast from "@/components/Toast";
import { api } from "@/lib/api";
import ServiceForm from "./ServiceForm";

interface Service {
  id: string;
  name: string;
  description?: string;
  basePrice: number;
  isActive: boolean;
}

const fetcher = (url: string) => api<Service[]>(url);

export default function ServicesPage() {
  const { data, isLoading, mutate } = useSWR("/services", fetcher);
  const [isOpen, setIsOpen] = useState(false);
  const [editing, setEditing] = useState<Service | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const closeModal = () => {
    setIsOpen(false);
    setEditing(null);
  };
  const openForNew = () => {
    setEditing(null);
    setIsOpen(true);
  };
  const openForEdit = (s: Service) => {
    setEditing(s);
    setIsOpen(true);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Remover serviço?")) return;
    try {
      await api(`/services/${id}`, { method: "DELETE" });
      setToast("Serviço removido");
      mutate();
    } catch {
      setToast("Erro ao remover serviço");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  const handleSuccess = () => {
    mutate();
    setToast("Serviço salvo");
    setTimeout(() => setToast(null), 3000);
    closeModal();
  };

  return (
    <ProtectedRoute roles={["admin"]}>
      <div className="p-4">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">Serviços</h1>
          <button onClick={openForNew} className="bg-primary text-background px-4 py-2">
            Novo Serviço
          </button>
        </div>
        {isLoading ? (
          <div className="flex justify-center p-4">
            <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-card">
              <thead className="bg-card sticky top-0">
                <tr>
                  <th className="px-4 py-2 text-left">Nome</th>
                  <th className="px-4 py-2 text-left">Descrição</th>
                  <th className="px-4 py-2 text-left">Preço</th>
                  <th className="px-4 py-2 text-left">Ativo</th>
                  <th className="px-4 py-2 text-left">Ações</th>
                </tr>
              </thead>
              <tbody>
                {data?.map((s, idx) => (
                  <tr key={s.id} className={idx % 2 ? "bg-card" : "bg-card"}>
                    <td className="px-4 py-2">{s.name}</td>
                    <td className="px-4 py-2">{s.description}</td>
                    <td className="px-4 py-2">
                      <Money value={s.basePrice} />
                    </td>
                    <td className="px-4 py-2">{s.isActive ? "Sim" : "Não"}</td>
                    <td className="px-4 py-2 space-x-2">
                      <button aria-label="Editar" onClick={() => openForEdit(s)}>
                        ✏️
                      </button>
                      <button aria-label="Remover" onClick={() => handleDelete(s.id)}>
                        🗑
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <Transition appear show={isOpen} as="div">
          <Dialog as="div" className="fixed inset-0 z-10 overflow-y-auto" onClose={closeModal}>
            <div className="min-h-screen px-4 text-center">
              <Transition.Child
                as="div"
                className="fixed inset-0 bg-black/50"
                enter="ease-out duration-300"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              />
              <span className="inline-block h-screen align-middle" aria-hidden="true">
                &#8203;
              </span>
              <Transition.Child
                as="div"
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
                className="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-card shadow-xl"
              >
                <ServiceForm service={editing || undefined} onClose={closeModal} onSuccess={handleSuccess} />
              </Transition.Child>
            </div>
          </Dialog>
        </Transition>
        <Toast message={toast} />
      </div>
    </ProtectedRoute>
  );
}
