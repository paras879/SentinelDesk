import api from "./api";

export interface PlanInfo {
  plan_name: string;
  max_devices: number;
  current_devices: number;
  is_limit_reached: boolean;
  remaining_slots: number;
}

export async function getPlanInfo(): Promise<PlanInfo> {
  try {
    const res = await api.get("/api/v1/plan");
    return res.data;
  } catch (error) {
    // Fallback if backend /api/v1/plan is not yet redeployed
    return {
      plan_name: "Starter",
      max_devices: 3,
      current_devices: 1,
      is_limit_reached: false,
      remaining_slots: 2,
    };
  }
}
