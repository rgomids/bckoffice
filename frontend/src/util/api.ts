const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function apiFetch<T = unknown>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  } as Record<string, string>;
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  if (!res.ok) throw new Error(await res.text());
  return (await res.json()) as T;
}

export async function apiPresign(
  path: string,
  filename: string,
): Promise<string> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers: Record<string, string> = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}${path}?filename=${encodeURIComponent(filename)}`, {
    method: "PUT",
    headers,
  });
  if (!res.ok) throw new Error(await res.text());
  const data = (await res.json()) as { url: string };
  return data.url;
}

export { apiFetch as api };
