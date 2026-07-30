import api from "./api";
import type { Device, DeviceSystemInfo, SoftwareItem, ProcessItem, WindowsServiceItem } from "@/types";

export async function getDevices(locationType?: string): Promise<{ count: number; devices: Device[] }> {
  const params = locationType ? { location_type: locationType } : {};
  const res = await api.get("/api/v1/devices", { params });
  return res.data;
}

export async function updateDeviceLocation(id: string, locationType: string): Promise<void> {
  await api.put(`/api/v1/devices/${id}/location`, { location_type: locationType });
}

export async function getDevice(id: string): Promise<Device> {
  const res = await api.get(`/api/v1/devices/${id}`);
  return res.data;
}

export async function getDeviceSystemInfo(id: string): Promise<DeviceSystemInfo> {
  const res = await api.get(`/api/v1/devices/${id}/system-info`);
  return res.data;
}

export async function getDeviceSoftware(id: string): Promise<{ count: number; software: SoftwareItem[] }> {
  const res = await api.get(`/api/v1/devices/${id}/software`);
  return res.data;
}

export async function getDeviceProcesses(id: string): Promise<{ count: number; processes: ProcessItem[] }> {
  const res = await api.get(`/api/v1/devices/${id}/processes`);
  return res.data;
}

export async function getDeviceServices(id: string): Promise<{ count: number; services: WindowsServiceItem[] }> {
  const res = await api.get(`/api/v1/devices/${id}/services`);
  return res.data;
}

export async function startService(deviceId: string, serviceName: string): Promise<void> {
  await api.post(`/api/v1/devices/${deviceId}/services/start`, { service_name: serviceName });
}

export async function stopService(deviceId: string, serviceName: string): Promise<void> {
  await api.post(`/api/v1/devices/${deviceId}/services/stop`, { service_name: serviceName });
}

export async function restartService(deviceId: string, serviceName: string): Promise<void> {
  await api.post(`/api/v1/devices/${deviceId}/services/restart`, { service_name: serviceName });
}

export async function deleteDevice(deviceId: string): Promise<void> {
  await api.delete(`/api/v1/devices/${deviceId}`);
}
