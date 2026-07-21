import api from "./api";
import type { LoginRequest, LoginResponse } from "@/types";

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await api.post<LoginResponse>("/api/v1/auth/login", data);
  return res.data;
}

export async function updateProfile(data: { name?: string; email?: string; password?: string }) {
  const res = await api.put("/api/v1/admin/profile", data);
  return res.data;
}
