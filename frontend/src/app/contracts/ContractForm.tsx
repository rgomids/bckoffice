"use client";
import { FormEvent, useEffect, useState } from "react";
import Attachments from "./Attachments";
import { api } from "@/lib/api";

interface ContractFormProps {
  contract?: Contract;
  onClose: () => void;
  onSuccess: () => void;
}

export interface Contract {
  id: string;
  customer?: { id: string; trade_name: string };
  service?: { id: string; name: string };
  promoter?: { id: string; name: string };
  value_total: number;
  start_date: string;
  end_date?: string;
  status: string;
}

interface Customer { id: string; trade_name: string; }
interface Service { id: string; name: string; }
interface Promoter { id: string; name: string; }

export default function ContractForm({ contract, onClose, onSuccess }: ContractFormProps) {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [promoters, setPromoters] = useState<Promoter[]>([]);

  const [customerInput, setCustomerInput] = useState(contract?.customer?.trade_name || "");
  const [customerID, setCustomerID] = useState(contract?.customer?.id || "");
  const [serviceID, setServiceID] = useState(contract?.service?.id || "");
  const [promoterID, setPromoterID] = useState(contract?.promoter?.id || "");
  const [valueTotal, setValueTotal] = useState(contract?.value_total ?? 0);
  const [startDate, setStartDate] = useState(contract?.start_date || "");
  const [endDate, setEndDate] = useState(contract?.end_date || "");
  const [status, setStatus] = useState(contract?.status || "active");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    (async () => {
      setCustomers(await api<Customer[]>("/customers"));
      setServices(await api<Service[]>("/services?active=true"));
      setPromoters(await api<Promoter[]>("/promoters"));
    })();
  }, []);

  useEffect(() => {
    const found = customers.find((c) => c.trade_name === customerInput);
    if (found) setCustomerID(found.id);
  }, [customerInput, customers]);

  const validate = () => {
    if (!customerID) return "Cliente obrigatório";
    if (!serviceID) return "Serviço obrigatório";
    if (valueTotal < 0) return "Valor inválido";
    if (startDate && endDate && new Date(startDate) > new Date(endDate))
      return "Datas incoerentes";
    return "";
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const err = validate();
    if (err) {
      setError(err);
      return;
    }
    setLoading(true);
    setError("");
    const payload = {
      customer_id: customerID,
      service_id: serviceID,
      promoter_id: promoterID || undefined,
      value_total: valueTotal,
      start_date: startDate,
      end_date: endDate,
      status,
    };
    try {
      if (contract) {
        await api(`/contracts/${contract.id}`, {
          method: "PUT",
          body: JSON.stringify(payload),
        });
      } else {
        await api("/contracts", { method: "POST", body: JSON.stringify(payload) });
      }
      onSuccess();
    } catch {
      setError("Erro ao salvar");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        <div>
          <input
            list="cust-list"
            className="border p-2 w-full"
            placeholder="Cliente"
            value={customerInput}
            onChange={(e) => setCustomerInput(e.target.value)}
            required
          />
          <datalist id="cust-list">
            {customers.map((c) => (
              <option key={c.id} value={c.trade_name} />
            ))}
          </datalist>
        </div>
        <select
          className="border p-2 w-full"
          value={serviceID}
          onChange={(e) => setServiceID(e.target.value)}
          required
        >
          <option value="" disabled>
            Serviço...
          </option>
          {services.map((s) => (
            <option key={s.id} value={s.id}>
              {s.name}
            </option>
          ))}
        </select>
        <select
          className="border p-2 w-full"
          value={promoterID}
          onChange={(e) => setPromoterID(e.target.value)}
        >
          <option value="">Promotor...</option>
          {promoters.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        <input
          type="number"
          min="0"
          step="0.01"
          className="border p-2 w-full"
          placeholder="Valor"
          value={valueTotal}
          onChange={(e) => setValueTotal(parseFloat(e.target.value))}
          required
        />
        <input
          type="date"
          className="border p-2 w-full"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
          required
        />
        <input
          type="date"
          className="border p-2 w-full"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
        />
        <select
          className="border p-2 w-full"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        >
          <option value="active">active</option>
          <option value="suspended">suspended</option>
          <option value="closed">closed</option>
          <option value="cancelled">cancelled</option>
        </select>
      </div>
      {contract && <Attachments contractId={contract.id} />}
      {error && <p className="text-red-500 text-sm">{error}</p>}
      <div className="flex gap-2 justify-end">
        <button type="button" onClick={onClose} className="border px-4 py-2">
          Cancelar
        </button>
        <button
          type="submit"
          className="bg-primary text-background px-4 py-2"
          disabled={loading}
        >
          {loading ? (
            <div className="h-4 w-4 border-2 border-card border-t-transparent rounded-full animate-spin" />
          ) : (
            "Salvar"
          )}
        </button>
      </div>
    </form>
  );
}
