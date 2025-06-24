"use client";
import { useState } from "react";
import useSWR from "swr";
import { DndContext, closestCenter, DragEndEvent } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import ProtectedRoute from "@/components/ProtectedRoute";
import LeadCard from "./LeadCard";
import NewLeadModal from "./NewLeadModal";
import { api } from "@/lib/api";

interface Lead {
  id: string;
  customer: { trade_name: string };
  service: { name: string };
  createdAt: string;
}

const fetcher = (url: string) => api<Lead[]>(url);

export default function LeadsPage() {
  const { data: leadsLead, mutate: mutLead } = useSWR<Lead[]>("/leads?status=lead", fetcher);
  const { data: leadsQualified, mutate: mutQualified } = useSWR<Lead[]>("/leads?status=qualified", fetcher);
  const { data: leadsProposal, mutate: mutProposal } = useSWR<Lead[]>("/leads?status=proposal", fetcher);
  const { data: leadsContract, mutate: mutContract } = useSWR<Lead[]>("/leads?status=contract", fetcher);

  const [open, setOpen] = useState(false);

  const refresh = () => {
    mutLead();
    mutQualified();
    mutProposal();
    mutContract();
  };

  const onDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const targetStatus = over.id as string;
    await api(`/leads/${active.id}/status`, { method: "PUT", body: JSON.stringify({ status: targetStatus }) });
    refresh();
  };

  return (
    <ProtectedRoute roles={["admin", "promoter"]}>
      <div className="p-4">
        <div className="flex justify-between mb-4">
          <h1 className="text-xl font-bold">Leads</h1>
          <button onClick={() => setOpen(true)} className="bg-primary text-background px-4 py-2">Novo Lead</button>
        </div>
        <DndContext onDragEnd={onDragEnd} collisionDetection={closestCenter}>
          <div className="grid grid-cols-4 gap-4">
            <SortableContext items={leadsLead?.map((l) => l.id) || []} strategy={verticalListSortingStrategy} id="lead">
              <div id="lead" className="bg-card p-2 min-h-[200px]">
                <h2 className="font-semibold mb-2">Lead</h2>
                {leadsLead?.map((l) => (
                  <LeadCard key={l.id} id={l.id} customer={l.customer} service={l.service} createdAt={l.createdAt} />
                ))}
              </div>
            </SortableContext>
            <SortableContext items={leadsQualified?.map((l) => l.id) || []} strategy={verticalListSortingStrategy} id="qualified">
              <div id="qualified" className="bg-card p-2 min-h-[200px]">
                <h2 className="font-semibold mb-2">Qualified</h2>
                {leadsQualified?.map((l) => (
                  <LeadCard key={l.id} id={l.id} customer={l.customer} service={l.service} createdAt={l.createdAt} />
                ))}
              </div>
            </SortableContext>
            <SortableContext items={leadsProposal?.map((l) => l.id) || []} strategy={verticalListSortingStrategy} id="proposal">
              <div id="proposal" className="bg-card p-2 min-h-[200px]">
                <h2 className="font-semibold mb-2">Proposal</h2>
                {leadsProposal?.map((l) => (
                  <LeadCard key={l.id} id={l.id} customer={l.customer} service={l.service} createdAt={l.createdAt} />
                ))}
              </div>
            </SortableContext>
            <SortableContext items={leadsContract?.map((l) => l.id) || []} strategy={verticalListSortingStrategy} id="contract">
              <div id="contract" className="bg-card p-2 min-h-[200px]">
                <h2 className="font-semibold mb-2">Contract</h2>
                {leadsContract?.map((l) => (
                  <LeadCard key={l.id} id={l.id} customer={l.customer} service={l.service} createdAt={l.createdAt} />
                ))}
              </div>
            </SortableContext>
          </div>
        </DndContext>
        <NewLeadModal isOpen={open} onClose={() => setOpen(false)} onSuccess={refresh} />
      </div>
    </ProtectedRoute>
  );
}
