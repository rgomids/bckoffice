"use client";

import { FormEvent, useState } from "react";
import { api } from "@/lib/api";

interface Service {
  id: string;
  name: string;
  description?: string;
  basePrice: number;
  isActive: boolean;
}

interface ServiceFormProps {
  service?: Service;
  onClose: () => void;
  onSuccess: () => void;
}

export default function ServiceForm({ service, onClose, onSuccess }: ServiceFormProps) {
  const [name, setName] = useState(service?.name || "");
  const [description, setDescription] = useState(service?.description || "");
  const [price, setPrice] = useState(service?.basePrice ?? 0);
  const [active, setActive] = useState(service?.isActive ?? true);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const validate = () => {
    if (!name) return "Nome é obrigatório";
    if (price < 0) return "Preço deve ser maior ou igual a 0";
    return "";
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const err = validate();
    if (err) {
      setError(err);
      return;
    }
    setError("");
    setLoading(true);
    const payload = {
      name,
      description,
      base_price: price,
      is_active: active,
    };
    try {
      if (service) {
        await api(`/services/${service.id}`, {
          method: "PUT",
          body: JSON.stringify(payload),
        });
      } else {
        await api("/services", {
          method: "POST",
          body: JSON.stringify(payload),
        });
      }
      onSuccess();
    } catch {
      setError("Erro ao salvar serviço");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-2">
      <input
        className="border p-2"
        placeholder="Nome"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
      <textarea
        className="border p-2"
        placeholder="Descrição"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <input
        type="number"
        min="0"
        step="0.01"
        className="border p-2"
        placeholder="Preço"
        value={price}
        onChange={(e) => setPrice(parseFloat(e.target.value))}
        required
      />
      <label className="flex items-center gap-2">
        <input type="checkbox" checked={active} onChange={(e) => setActive(e.target.checked)} />
        Ativo
      </label>
      {error && <p className="text-red-500 text-sm">{error}</p>}
      <div className="flex gap-2 mt-2">
        <button type="button" onClick={onClose} className="border px-4 py-2">
          Cancelar
        </button>
        <button type="submit" className="bg-primary text-background px-4 py-2" disabled={loading}>
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
