"use client";

import { FormEvent, useState } from "react";
import { api } from "@/util/api";

interface CustomerFormProps {
  onClose: () => void;
  onSuccess: () => void;
  customer?: Customer;
}

interface Customer {
  id: string;
  legalName: string;
  tradeName: string;
  documentID: string;
  email: string;
  phone: string;
}

export default function CustomerForm({ onClose, onSuccess, customer }: CustomerFormProps) {
  const [legalName, setLegalName] = useState(customer?.legalName || "");
  const [tradeName, setTradeName] = useState(customer?.tradeName || "");
  const [documentId, setDocumentId] = useState(customer?.documentID || "");
  const [email, setEmail] = useState(customer?.email || "");
  const [phone, setPhone] = useState(customer?.phone || "");
  const [street, setStreet] = useState("");
  const [number, setNumber] = useState("");
  const [city, setCity] = useState("");
  const [state, setState] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [error, setError] = useState("");

  const validate = () => {
    if (!legalName || !documentId) return false;
    if (email && !/^\S+@\S+\.\S+$/.test(email)) return false;
    return true;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validate()) {
      setError("Dados inválidos");
      return;
    }

    const payload = {
      legal_name: legalName,
      trade_name: tradeName,
      document_id: documentId,
      email,
      phone,
      addresses: [
        {
          address_type: "main",
          street,
          number,
          city,
          state,
          postal_code: postalCode,
        },
      ],
    };

    const method = customer ? "PUT" : "POST";
    const url = customer ? `/customers/${customer.id}` : "/customers";

    await api(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    onSuccess();
    onClose();
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-2 p-4">
      <input
        className="border p-2"
        placeholder="Razão Social"
        value={legalName}
        onChange={(e) => setLegalName(e.target.value)}
        required
      />
      <input
        className="border p-2"
        placeholder="Nome Fantasia"
        value={tradeName}
        onChange={(e) => setTradeName(e.target.value)}
      />
      <input
        className="border p-2"
        placeholder="Documento"
        value={documentId}
        onChange={(e) => setDocumentId(e.target.value)}
        required
      />
      <input
        className="border p-2"
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        className="border p-2"
        placeholder="Telefone"
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
      />

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        <input
          className="border p-2"
          placeholder="Rua"
          value={street}
          onChange={(e) => setStreet(e.target.value)}
        />
        <input
          className="border p-2"
          placeholder="Número"
          value={number}
          onChange={(e) => setNumber(e.target.value)}
        />
        <input
          className="border p-2"
          placeholder="Cidade"
          value={city}
          onChange={(e) => setCity(e.target.value)}
        />
        <input
          className="border p-2"
          placeholder="Estado"
          value={state}
          onChange={(e) => setState(e.target.value)}
        />
        <input
          className="border p-2"
          placeholder="CEP"
          value={postalCode}
          onChange={(e) => setPostalCode(e.target.value)}
        />
      </div>

      {error && <p className="text-red-500 text-sm">{error}</p>}

      <div className="flex gap-2 mt-2">
        <button type="button" onClick={onClose} className="border px-4 py-2">
          Cancelar
        </button>
        <button type="submit" className="bg-blue-500 text-white px-4 py-2">
          Salvar
        </button>
      </div>
    </form>
  );
}

