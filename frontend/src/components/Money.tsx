"use client";

interface MoneyProps {
  value: number;
}

export default function Money({ value }: MoneyProps) {
  const formatted = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(value);
  return <span>{formatted}</span>;
}
