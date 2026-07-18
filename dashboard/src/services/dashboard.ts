import api from "./api";
import type { DashboardSummary } from "@/types";

export async function getDashboardSummary(): Promise<DashboardSummary> {
  const res = await api.get<DashboardSummary>("/api/v1/dashboard/summary");
  return res.data;
}

export interface MyIPResponse {
  ip: string;
  network_group_id: string;
}

export async function getMyIP(): Promise<MyIPResponse> {
  const res = await api.get<MyIPResponse>("/api/v1/dashboard/my-ip");
  return res.data;
}
