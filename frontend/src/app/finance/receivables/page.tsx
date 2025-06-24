import { Suspense } from "react";
import ProtectedRoute from "@/components/ProtectedRoute";
import ReceivablesClient from "./receivables-client";

export default function ReceivablesPage() {
  return (
    <ProtectedRoute roles={["finance", "admin"]}>
      <Suspense fallback={<div className="p-4">Carregando...</div>}>
        <ReceivablesClient />
      </Suspense>
    </ProtectedRoute>
  );
}
