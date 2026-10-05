"use client";
import { CSS } from "@dnd-kit/utilities";
import { useSortable } from "@dnd-kit/sortable";

interface LeadCardProps {
  id: string;
  customer: { trade_name: string };
  service: { name: string };
  createdAt: string;
}

export default function LeadCard({ id, customer, service, createdAt }: LeadCardProps) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id });
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };
  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className="border border-card p-2 bg-card mb-2"
    >
      <div className="font-semibold">{customer.trade_name}</div>
      <div className="text-sm">{service.name}</div>
      <div className="text-xs text-foreground">{new Date(createdAt).toLocaleDateString()}</div>
    </div>
  );
}
