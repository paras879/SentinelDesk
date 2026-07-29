"use client";

import { useEffect, useRef } from "react";
import { useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { getDevices } from "@/services/devices";
import { MonitorPlay } from "lucide-react";
import type { Device } from "@/types";
import { useAppNotifications } from "@/contexts/NotificationContext";

export function DeviceNotifier() {
  const previousDevicesRef = useRef<Set<string>>(new Set());
  const { addNotification } = useAppNotifications();

  const { data } = useQuery({
    queryKey: ["devices", "notifier"],
    queryFn: () => getDevices(),
    refetchInterval: 5000,
  });

  useEffect(() => {
    if (!data?.devices) return;

    const currentDevices = new Set(data.devices.map((d: Device) => d.ID));
    
    // Check for new devices
    if (previousDevicesRef.current.size > 0) {
      const newDevices = data.devices.filter((d: Device) => !previousDevicesRef.current.has(d.ID));
      
      newDevices.forEach((device: Device) => {
        addNotification({
          deviceId: device.ID,
          message: `New Agent Registered: ${device.DeviceName || device.Hostname || "Unknown"}`,
        });

        toast("New Device Registered!", {
          description: `${device.DeviceName || device.Hostname || "Unknown"} has joined the network.`,
          icon: <MonitorPlay className="h-4 w-4 text-emerald-500" />,
        });
      });
    }

    previousDevicesRef.current = currentDevices;
  }, [data?.devices]);

  return null; // Component does not render anything
}
