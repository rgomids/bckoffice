"use client";
import { useState, Fragment } from "react";
import useSWR from "swr";
import { Dialog, Transition } from "@headlessui/react";
import ProtectedRoute from "@/components/ProtectedRoute";
import Money from "@/components/Money";
import Toast from "@/components/Toast";
import { api } from "@/util/api";
import ContractForm, { Contract } from "./ContractForm";
import Attachments from "./Attachments";

const fetcher = (url: string) => api<Contract[]>(url);

export default function ContractsPage() {
  const { data, isLoading, mutate } = useSWR(
    "/contracts?status=active",
    fetcher,
  );
  const [isFormOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Contract | null>(null);
  const [attId, setAttId] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const closeForm = () => {
    setFormOpen(false);
    setEditing(null);
  };
  const openForNew = () => {
    setEditing(null);
    setFormOpen(true);
  };
  const openForEdit = (c: Contract) => {
    setEditing(c);
    setFormOpen(true);
  };
  const closeAtt = () => setAttId(null);

  const handleDelete = async (id: string) => {
    if (!confirm("Remover contrato?")) return;
    try {
      await api(`/contracts/${id}`, { method: "DELETE" });
      setToast("Contrato removido");
      mutate();
    } catch {
      setToast("Erro ao remover contrato");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  const handleSuccess = () => {
    mutate();
    setToast("Contrato salvo");
    setTimeout(() => setToast(null), 3000);
    closeForm();
  };

  return (
    <ProtectedRoute roles={["admin", "finance"]}>
      <div className="p-4">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">Contratos</h1>
          <button onClick={openForNew} className="bg-blue-500 text-white px-4 py-2">
            Novo Contrato
          </button>
        </div>
        {isLoading ? (
          <div className="flex justify-center p-4">
            <div className="h-5 w-5 border-2 border-gray-300 border-t-transparent rounded-full animate-spin" />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left">Cliente</th>
                  <th className="px-4 py-2 text-left">Serviço</th>
                  <th className="px-4 py-2 text-left">Valor</th>
                  <th className="px-4 py-2 text-left">Início</th>
                  <th className="px-4 py-2 text-left">Término</th>
                  <th className="px-4 py-2 text-left">Status</th>
                  <th className="px-4 py-2 text-left">Ações</th>
                </tr>
              </thead>
              <tbody>
                {data?.map((c, idx) => (
                  <tr key={c.id} className={idx % 2 ? "bg-gray-50" : "bg-white"}>
                    <td className="px-4 py-2">{c.customer?.trade_name}</td>
                    <td className="px-4 py-2">{c.service?.name}</td>
                    <td className="px-4 py-2">
                      <Money value={c.value_total} />
                    </td>
                    <td className="px-4 py-2">
                      {new Date(c.start_date).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-2">
                      {c.end_date ? new Date(c.end_date).toLocaleDateString() : ""}
                    </td>
                    <td className="px-4 py-2">{c.status}</td>
                    <td className="px-4 py-2 space-x-2">
                      <button onClick={() => openForEdit(c)}>✏️</button>
                      <button onClick={() => setAttId(c.id)}>📎</button>
                      <button onClick={() => handleDelete(c.id)}>🗑</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <Transition appear show={isFormOpen} as={Fragment}>
          <Dialog as="div" className="fixed inset-0 z-10 overflow-y-auto" onClose={closeForm}>
            <div className="min-h-screen px-4 text-center">
              <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
                <div className="fixed inset-0 bg-black/50" />
              </Transition.Child>
              <span className="inline-block h-screen align-middle" aria-hidden="true">
                &#8203;
              </span>
              <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
                <div className="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl">
                  <ContractForm contract={editing || undefined} onClose={closeForm} onSuccess={handleSuccess} />
                </div>
              </Transition.Child>
            </div>
          </Dialog>
        </Transition>
        <Transition appear show={!!attId} as={Fragment}>
          <Dialog as="div" className="fixed inset-0 z-10 overflow-y-auto" onClose={closeAtt}>
            <div className="min-h-screen px-4 text-center">
              <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
                <div className="fixed inset-0 bg-black/50" />
              </Transition.Child>
              <span className="inline-block h-screen align-middle" aria-hidden="true">
                &#8203;
              </span>
              <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
                <div className="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl">
                  {attId && <Attachments contractId={attId} />}
                </div>
              </Transition.Child>
            </div>
          </Dialog>
        </Transition>
        <Toast message={toast} />
      </div>
    </ProtectedRoute>
  );
}
