import { getToken } from "@/hooks/useAuth";

export async function api<T>(
  input: RequestInfo,
  init: RequestInit = {},
): Promise<T> {
  const token = getToken();
  const headers = new Headers(init.headers);
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  const response = await fetch(input, { ...init, headers });
  return response.json() as Promise<T>;
}
