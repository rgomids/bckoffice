"use client";

import { ReactNode, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getToken } from "@/hooks/useAuth";

interface Props {
  children: ReactNode;
  roles?: string[];
}

function parseJwt(token: string): any {
  try {
    return JSON.parse(atob(token.split(".")[1]));
  } catch {
    return null;
  }
}
export default function ProtectedRoute({ children, roles }: Props) {
  const router = useRouter();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
      return;
    }
    if (roles && roles.length > 0) {
      const payload = parseJwt(token);
      if (!payload || !roles.includes(payload.role)) {
        router.push("/login");
      }
    }
  }, [router, roles]);

  return <>{children}</>;
}
