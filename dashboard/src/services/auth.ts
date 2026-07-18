import api from "./api";
import type { LoginRequest, LoginResponse } from "@/types";

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await api.post<LoginResponse>("/api/v1/auth/login", data);
  return res.data;
}
