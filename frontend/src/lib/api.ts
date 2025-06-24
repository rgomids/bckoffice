const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

import { mutate } from "swr";

export interface ApiError {
  message: string;
  status: number;
}

export async function api<T = unknown>(
  path: string,
  opts: RequestInit = {},
): Promise<T> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers = {
    "Content-Type": "application/json",
    ...(opts.headers || {}),
  } as Record<string, string>;
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });

  if (res.status === 401 && typeof window !== "undefined") {
    window.location.href = "/login";
  }

  if (!res.ok) {
    const text = await res.text();
    throw { message: text || res.statusText, status: res.status } as ApiError;
  }

  if (opts.method && opts.method !== "GET") {
    mutate(path);
  }

  if (res.status === 204) {
    return null as unknown as T;
  }
  return (await res.json()) as T;
}

export async function apiPresign(path: string, filename: string): Promise<string> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers: Record<string, string> = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const res = await fetch(
    `${API_BASE}${path}?filename=${encodeURIComponent(filename)}`,
    {
      method: "PUT",
      headers,
    },
  );
  if (!res.ok) {
    const text = await res.text();
    throw { message: text || res.statusText, status: res.status } as ApiError;
  }
  const data = (await res.json()) as { url: string };
  return data.url;
}
