"use client";

import { useEffect, useState } from "react";
import { Dialog, Transition } from "@headlessui/react";
import ProtectedRoute from "@/components/ProtectedRoute";
import CustomerForm from "./CustomerForm";
import { api } from "@/lib/api";

interface Customer {
  id: string;
  legalName: string;
  tradeName: string;
  documentID: string;
  email: string;
  phone: string;
}

export default function CustomersPage() {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [editing, setEditing] = useState<Customer | null>(null);

  const fetchCustomers = async () => {
    const data = await api<Customer[]>("/customers");
    setCustomers(data);
  };

  useEffect(() => {
    fetchCustomers();
  }, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Remover cliente?")) return;
    await api(`/customers/${id}`, { method: "DELETE" });
    fetchCustomers();
  };

  const closeModal = () => {
    setIsOpen(false);
    setEditing(null);
  };

  const openForNew = () => {
    setEditing(null);
    setIsOpen(true);
  };

  const openForEdit = (c: Customer) => {
    setEditing(c);
    setIsOpen(true);
  };

  return (
    <ProtectedRoute roles={["admin", "finance"]}>
      <div className="p-4">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold">Clientes</h1>
          <button onClick={openForNew} className="bg-primary text-background px-4 py-2">
            Novo Cliente
          </button>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-card">
            <thead className="bg-card sticky top-0">
              <tr>
                <th className="px-4 py-2 text-left">Legal Name</th>
                <th className="px-4 py-2 text-left">Trade Name</th>
                <th className="px-4 py-2 text-left">Document</th>
                <th className="px-4 py-2 text-left">Email</th>
                <th className="px-4 py-2 text-left">Phone</th>
                <th className="px-4 py-2 text-left">Actions</th>
              </tr>
            </thead>
            <tbody>
              {customers.map((c, idx) => (
                <tr key={c.id} className={idx % 2 ? "bg-card" : "bg-card"}>
                  <td className="px-4 py-2">{c.legalName}</td>
                  <td className="px-4 py-2">{c.tradeName}</td>
                  <td className="px-4 py-2">{c.documentID}</td>
                  <td className="px-4 py-2">{c.email}</td>
                  <td className="px-4 py-2">{c.phone}</td>
                  <td className="px-4 py-2 space-x-2">
                    <button aria-label="Editar" onClick={() => openForEdit(c)}>
                      ✏️
                    </button>
                    <button aria-label="Remover" onClick={() => handleDelete(c.id)}>
                      🗑
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

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
                <CustomerForm
                  onClose={closeModal}
                  onSuccess={fetchCustomers}
                  customer={editing || undefined}
                />
              </Transition.Child>
            </div>
          </Dialog>
        </Transition>
      </div>
    </ProtectedRoute>
  );
}

