"use client";
import { Dialog, Transition } from "@headlessui/react";
import { Fragment, FormEvent, useState, useEffect } from "react";
import { api } from "@/util/api";
import { getToken } from "@/hooks/useAuth";

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

interface Customer { id: string; trade_name: string; }
interface Service { id: string; name: string; }

export default function NewLeadModal({ isOpen, onClose, onSuccess }: ModalProps) {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [customerID, setCustomerID] = useState("");
  const [serviceID, setServiceID] = useState("");
  const [notes, setNotes] = useState("");

  useEffect(() => {
    (async () => {
      setCustomers(await api<Customer[]>("/customers"));
      setServices(await api<Service[]>("/services"));
    })();
  }, []);

  interface LeadPayload {
    customer_id: string;
    service_id: string;
    notes: string;
    promoter_id?: string;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const payload: LeadPayload = {
      customer_id: customerID,
      service_id: serviceID,
      notes,
    };
    const token = getToken();
    if (token) {
      const payloadJwt = JSON.parse(atob(token.split(".")[1]));
      if (payloadJwt.role === "promoter") payload.promoter_id = payloadJwt.sub;
    }
    await api("/leads", { method: "POST", body: JSON.stringify(payload) });
    onSuccess();
    onClose();
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="fixed inset-0 z-10 overflow-y-auto" onClose={onClose}>
        <div className="min-h-screen px-4 text-center">
          <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
            <div className="fixed inset-0 bg-black/50" />
          </Transition.Child>
          <span className="inline-block h-screen align-middle" aria-hidden="true">&#8203;</span>
          <Transition.Child as={Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
            <div className="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl">
              <form onSubmit={handleSubmit} className="flex flex-col gap-2">
                <select className="border p-2" value={customerID} onChange={(e) => setCustomerID(e.target.value)} required>
                  <option value="" disabled>Cliente...</option>
                  {customers.map((c) => (
                    <option key={c.id} value={c.id}>{c.trade_name}</option>
                  ))}
                </select>
                <select className="border p-2" value={serviceID} onChange={(e) => setServiceID(e.target.value)} required>
                  <option value="" disabled>Serviço...</option>
                  {services.map((s) => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
                <textarea className="border p-2" placeholder="Notas" value={notes} onChange={(e) => setNotes(e.target.value)} />
                <div className="flex gap-2 mt-2">
                  <button type="button" onClick={onClose} className="border px-4 py-2">Cancelar</button>
                  <button type="submit" className="bg-blue-500 text-white px-4 py-2">Salvar</button>
                </div>
              </form>
            </div>
          </Transition.Child>
        </div>
      </Dialog>
    </Transition>
  );
}
