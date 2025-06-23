"use client";

import { ReactNode, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getToken } from "@/hooks/useAuth";

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const router = useRouter();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
    }
  }, [router]);

  return <>{children}</>;
}
