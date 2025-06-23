"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

import { apiFetch } from "@/util/api";

export const login = async (email: string, password: string) => {
  const data = await apiFetch("/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  localStorage.setItem("token", data.token);
};

export const logout = () => {
  localStorage.removeItem("token");
};

export const getToken = (): string | null => {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
};

interface AuthContextType {
  user: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<string | null>(null);

  useEffect(() => {
    if (getToken()) {
      setUser("authenticated");
    }
  }, []);

  const handleLogin = async (email: string, password: string) => {
    await login(email, password);
    setUser(email);
  };

  const handleLogout = () => {
    logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{ user, login: handleLogin, logout: handleLogout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
