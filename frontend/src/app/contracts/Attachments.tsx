"use client";
import { useRef, useState } from "react";
import useSWR from "swr";
import Toast from "@/components/Toast";
import { api, apiPresign } from "@/lib/api";

interface Attachment {
  id: string;
  filename: string;
  url: string;
}

const fetcher = (url: string) => api<Attachment[]>(url);

export default function Attachments({ contractId }: { contractId: string }) {
  const { data, isLoading, mutate } = useSWR(
    `/contracts/${contractId}/attachments`,
    fetcher,
  );
  const inputRef = useRef<HTMLInputElement>(null);
  const [toast, setToast] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);

  const selectFile = () => inputRef.current?.click();

  const onFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    e.target.value = "";
    try {
      setUploading(true);
      const url = await apiPresign(
        `/contracts/${contractId}/attachments/presign`,
        file.name,
      );
      await fetch(url, { method: "PUT", body: file });
      await api(`/contracts/${contractId}/attachments`, {
        method: "POST",
        body: JSON.stringify({ filename: file.name, url }),
      });
      setToast("Arquivo enviado");
      mutate();
    } catch {
      setToast("Erro ao enviar arquivo");
    } finally {
      setUploading(false);
      setTimeout(() => setToast(null), 3000);
    }
  };

  const remove = async (id: string) => {
    if (!confirm("Remover arquivo?")) return;
    try {
      await api(`/contracts/${contractId}/attachments/${id}`, { method: "DELETE" });
      setToast("Arquivo removido");
      mutate();
    } catch {
      setToast("Erro ao remover arquivo");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  return (
    <div className="space-y-2">
      {isLoading ? (
        <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
      ) : (
        <ul className="space-y-1">
          {data?.map((a) => (
            <li key={a.id} className="flex justify-between items-center">
              <a href={a.url} target="_blank" rel="noopener noreferrer" className="text-primary underline">
                {a.filename}
              </a>
              <button aria-label="Remover" onClick={() => remove(a.id)}>
                🗑
              </button>
            </li>
          ))}
        </ul>
      )}
      <input
        type="file"
        ref={inputRef}
        onChange={onFileChange}
        className="hidden"
      />
      <button
        type="button"
        onClick={selectFile}
        className="border px-2 py-1"
        disabled={uploading}
      >
        {uploading ? (
          <div className="h-4 w-4 border-2 border-card border-t-transparent rounded-full animate-spin" />
        ) : (
          "Adicionar Arquivo"
        )}
      </button>
      <Toast message={toast} />
    </div>
  );
}
