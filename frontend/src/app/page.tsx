"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { Transition } from "@headlessui/react";
import { motion } from "framer-motion";
import { api } from "@/lib/api";

interface Stats {
  totalCustomers: number;
  activeLeads: number;
  monthlyRevenue: number;
}

export default function Home() {
  const [stats, setStats] = useState<Stats | null>(null);

  useEffect(() => {
    api<Stats>("/stats").then(setStats).catch(() => {});
  }, []);

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 p-4">
      <section className="md:col-span-2 space-y-2">
        <h1 className="text-lg font-semibold">Visão Geral</h1>
        {stats ? (
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="bg-card border border-card rounded p-4">
              <p className="text-sm">Total Clientes</p>
              <p className="text-2xl font-bold">{stats.totalCustomers}</p>
            </div>
            <div className="bg-card border border-card rounded p-4">
              <p className="text-sm">Leads Ativos</p>
              <p className="text-2xl font-bold">{stats.activeLeads}</p>
            </div>
            <div className="bg-card border border-card rounded p-4">
              <p className="text-sm">$$ Receitas mês</p>
              <p className="text-2xl font-bold">
                {new Intl.NumberFormat("pt-BR", {
                  style: "currency",
                  currency: "BRL",
                }).format(stats.monthlyRevenue)}
              </p>
            </div>
          </div>
        ) : (
          <Transition
            show={true}
            as={motion.div}
            enter="transition-opacity"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            className="grid grid-cols-3 gap-4"
          >
            {Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="h-20 bg-card border border-card rounded animate-pulse"
              />
            ))}
          </Transition>
        )}
      </section>
      <section className="space-y-2">
        <h2 className="text-lg font-semibold">Atalhos</h2>
        <div className="grid gap-2">
          <Link
            href="/customers"
            className="bg-card border border-card rounded p-2 text-center"
          >
            Clientes
          </Link>
          <Link
            href="/leads"
            className="bg-card border border-card rounded p-2 text-center"
          >
            Leads
          </Link>
          <Link
            href="/finance/receivables"
            className="bg-card border border-card rounded p-2 text-center"
          >
            Financeiro
          </Link>
          <Link
            href="/services"
            className="bg-card border border-card rounded p-2 text-center"
          >
            Serviços
          </Link>
        </div>
      </section>
      <section className="md:col-span-3 mt-4 p-6 rounded text-center text-background bg-gradient-to-r from-primary to-secondary">
        <h2 className="text-3xl font-bold">RGPS</h2>
      </section>
    </div>
  );
}
