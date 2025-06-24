"use client";

import { Fragment, FormEvent, useEffect, useState } from "react";
import useSWR from "swr";
import { Dialog, Transition } from "@headlessui/react";
import ProtectedRoute from "@/components/ProtectedRoute";
import Money from "@/components/Money";
import Toast from "@/components/Toast";
import { api } from "@/lib/api";
import { getToken } from "@/hooks/useAuth";

interface Lead {
  id: string;
  status: string;
  customer: { trade_name: string };
  service: { name: string };
}

interface Commission {
  id: string;
  amount: number;
  approved: boolean;
}

interface Contract {
  id: string;
  customer: { trade_name: string };
  service: { name: string };
  value_total: number;
  start_date: string;
}

interface BankAccount {
  pix: string;
  bank: string;
  agency: string;
  account: string;
}

interface Promoter {
  fullName: string;
  bankAccount: BankAccount;
}

const fetcherLeads = (url: string) => api<Lead[]>(url);
const fetcherComms = (url: string) => api<Commission[]>(url);
const fetcherPromoter = (url: string) => api<Promoter>(url);
const fetcherContracts = (url: string) => api<Contract[]>(url);

function parseJwt(token: string): { sub?: string } | null {
  try {
    return JSON.parse(atob(token.split(".")[1]));
  } catch {
    return null;
  }
}

export default function PromoterDashboardPage() {
  const token = getToken();
  const userId = token ? parseJwt(token)?.sub : null;

  const { data: leads, isLoading: loadingLeads } = useSWR<Lead[]>(
    userId
      ? `/leads?promoter_id=${userId}&status=lead,qualified,proposal,contract`
      : null,
    fetcherLeads,
  );
  const { data: commissions, isLoading: loadingComms } = useSWR<Commission[]>(
    userId ? `/commissions?promoter_id=${userId}` : null,
    fetcherComms,
  );
  const { data: promoter, mutate: mutPromoter, isLoading: loadingProm } =
    useSWR<Promoter>(userId ? `/promoters/${userId}` : null, fetcherPromoter);
  const { data: contracts, isLoading: loadingContracts } = useSWR<Contract[]>(
    userId ? `/contracts?promoter_id=${userId}` : null,
    fetcherContracts,
  );

  const leadsOpen = leads?.filter((l) => l.status === "lead").length || 0;
  const leadsProposal =
    leads?.filter((l) => l.status === "proposal").length || 0;
  const pendingComms =
    commissions?.filter((c) => !c.approved).reduce((s, c) => s + c.amount, 0) ||
    0;

  const [isModalOpen, setModalOpen] = useState(false);
  const [pix, setPix] = useState("");
  const [bank, setBank] = useState("");
  const [agency, setAgency] = useState("");
  const [account, setAccount] = useState("");
  const [toast, setToast] = useState<string | null>(null);

  useEffect(() => {
    if (promoter?.bankAccount) {
      setPix(promoter.bankAccount.pix || "");
      setBank(promoter.bankAccount.bank || "");
      setAgency(promoter.bankAccount.agency || "");
      setAccount(promoter.bankAccount.account || "");
    }
  }, [promoter]);

  const closeModal = () => setModalOpen(false);

  const submitBank = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await api(`/promoters/${userId}`, {
        method: "PUT",
        body: JSON.stringify({
          full_name: promoter?.fullName || "",
          bank_account: { pix, bank, agency, account },
        }),
      });
      setToast("Dados atualizados");
      mutPromoter();
      closeModal();
    } catch {
      setToast("Erro ao salvar");
    } finally {
      setTimeout(() => setToast(null), 3000);
    }
  };

  const contractsList = (contracts || [])
    .sort((a, b) =>
      a.start_date < b.start_date ? 1 : a.start_date > b.start_date ? -1 : 0,
    )
    .slice(0, 10);

  return (
    <ProtectedRoute roles={["promoter"]}>
      <div className="p-4 space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-card p-4 shadow">
            <h2 className="font-semibold mb-2">Leads</h2>
            {loadingLeads ? (
              <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
            ) : (
              <div className="flex justify-between">
                <div>
                  <p className="text-sm text-foreground">Abertos</p>
                  <p className="text-xl font-bold">{leadsOpen}</p>
                </div>
                <div>
                  <p className="text-sm text-foreground">Proposta</p>
                  <p className="text-xl font-bold">{leadsProposal}</p>
                </div>
              </div>
            )}
          </div>
          <div className="bg-card p-4 shadow">
            <h2 className="font-semibold mb-2">Comissões Pendentes</h2>
            {loadingComms ? (
              <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
            ) : (
              <div className="text-xl font-bold">
                <Money value={pendingComms} />
              </div>
            )}
          </div>
        </div>
        <div className="bg-card p-4 shadow">
          <h2 className="font-semibold mb-2">Contratos Fechados</h2>
          {loadingContracts ? (
            <div className="flex justify-center p-4">
              <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-card">
                <thead className="bg-card">
                  <tr>
                    <th className="px-4 py-2 text-left">Cliente</th>
                    <th className="px-4 py-2 text-left">Serviço</th>
                    <th className="px-4 py-2 text-left">Valor</th>
                    <th className="px-4 py-2 text-left">Data Início</th>
                  </tr>
                </thead>
                <tbody>
                  {contractsList.map((c, idx) => (
                    <tr
                      key={c.id}
                      className={idx % 2 ? "bg-card" : "bg-card"}
                    >
                      <td className="px-4 py-2">{c.customer?.trade_name}</td>
                      <td className="px-4 py-2">{c.service?.name}</td>
                      <td className="px-4 py-2">
                        <Money value={c.value_total} />
                      </td>
                      <td className="px-4 py-2">
                        {new Date(c.start_date).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
        <div className="bg-card p-4 shadow">
          <div className="flex justify-between items-center mb-2">
            <h2 className="font-semibold">Meus Dados Bancários</h2>
            <button className="text-primary" onClick={() => setModalOpen(true)}>
              Editar
            </button>
          </div>
          {loadingProm ? (
            <div className="h-5 w-5 border-2 border-card border-t-transparent rounded-full animate-spin" />
          ) : (
            <div className="space-y-1 text-sm">
              <p>PIX: {promoter?.bankAccount?.pix}</p>
              <p>Banco: {promoter?.bankAccount?.bank}</p>
              <p>Agência: {promoter?.bankAccount?.agency}</p>
              <p>Conta: {promoter?.bankAccount?.account}</p>
            </div>
          )}
        </div>
        <Transition appear show={isModalOpen} as={Fragment}>
          <Dialog
            as="div"
            className="fixed inset-0 z-10 overflow-y-auto"
            onClose={closeModal}
          >
            <div className="min-h-screen px-4 text-center">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <div className="fixed inset-0 bg-black/50" />
              </Transition.Child>
              <span className="inline-block h-screen align-middle" aria-hidden="true">
                &#8203;
              </span>
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <div className="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-card shadow-xl">
                  <form onSubmit={submitBank} className="flex flex-col gap-2">
                    <input
                      className="border p-2"
                      placeholder="PIX"
                      value={pix}
                      onChange={(e) => setPix(e.target.value)}
                    />
                    <input
                      className="border p-2"
                      placeholder="Banco"
                      value={bank}
                      onChange={(e) => setBank(e.target.value)}
                    />
                    <input
                      className="border p-2"
                      placeholder="Agência"
                      value={agency}
                      onChange={(e) => setAgency(e.target.value)}
                    />
                    <input
                      className="border p-2"
                      placeholder="Conta"
                      value={account}
                      onChange={(e) => setAccount(e.target.value)}
                    />
                    <div className="flex gap-2 mt-2">
                      <button type="button" onClick={closeModal} className="border px-4 py-2">
                        Cancelar
                      </button>
                      <button type="submit" className="bg-primary text-background px-4 py-2">
                        Salvar
                      </button>
                    </div>
                  </form>
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
