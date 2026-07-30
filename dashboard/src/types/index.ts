export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  token: string;
  user: {
    id: string;
    name: string;
    email: string;
    role: string;
  };
}

export interface DashboardSummary {
  totalDevices: number;
  onlineDevices: number;
  offlineDevices: number;
}

export interface Device {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  DeviceName: string;
  Hostname: string;
  Username: string;
  OS: string;
  OSVersion: string;
  IPAddress: string;
  MACAddress: string;
  ConnectedSubnet?: string;
  NetworkAdapters?: any;
  DefaultGateway?: string;
  NetworkGroupID?: string;
  LocationType?: string;
  Status: "online" | "offline";
  LastSeen?: string | null;
}

export interface DeviceSystemInfo {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  device_id: string;
  cpu_name: string;
  cpu_usage: number;
  total_ram: number;
  used_ram: number;
  free_ram: number;
  total_disk: number;
  used_disk: number;
  free_disk: number;
  os_version: string;
  local_ip: string;
  public_ip: string;
  mac_address: string;
  uptime: number;
  last_boot: string;
  agent_version: string;
  Device?: Device;
}

export interface SoftwareItem {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  device_id: string;
  name: string;
  version: string;
  publisher: string;
  install_date: string;
  install_location: string;
  estimated_size: number;
}

export interface ProcessItem {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  device_id: string;
  pid: number;
  name: string;
  executable_path: string;
  cpu_usage: number;
  memory_bytes: number;
  memory_percent: number;
  thread_count: number;
  handle_count: number;
  start_time: string;
  username: string;
  status: string;
}

export interface WindowsServiceItem {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  device_id: string;
  service_name: string;
  display_name: string;
  status: string;
  start_type: string;
  executable_path: string;
  pid: number;
  service_account: string;
  description: string;
  can_stop: boolean;
  can_pause: boolean;
  accept_shutdown: boolean;
}
