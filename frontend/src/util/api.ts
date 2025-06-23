const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function apiFetch(
  path: string,
  options: RequestInit = {},
): Promise<any> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  } as Record<string, string>;
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export { apiFetch as api };
