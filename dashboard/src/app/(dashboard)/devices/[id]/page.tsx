"use client";

import { useState } from "react";
import { useParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { getDevice, getDeviceSystemInfo, getDeviceSoftware, getDeviceProcesses, getDeviceServices } from "@/services/devices";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { AlertTriangle, ArrowLeft, Monitor } from "lucide-react";
import { OverviewTab } from "@/features/devices/tabs/overview-tab";
import { SystemInfoTab } from "@/features/devices/tabs/system-info-tab";
import { SoftwareTab } from "@/features/devices/tabs/software-tab";
import { ProcessesTab } from "@/features/devices/tabs/processes-tab";
import { ServicesTab } from "@/features/devices/tabs/services-tab";
import { LiveViewTab } from "@/features/devices/tabs/live-view-tab";

type Tab = "overview" | "system-info" | "software" | "processes" | "services" | "live-view";

const tabs: { key: Tab; label: string }[] = [
  { key: "overview", label: "Overview" },
  { key: "live-view", label: "Live View" },
  { key: "system-info", label: "System Info" },
  { key: "software", label: "Software" },
  { key: "processes", label: "Processes" },
  { key: "services", label: "Services" },
];

export default function DeviceDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [activeTab, setActiveTab] = useState<Tab>("overview");

  const deviceQuery = useQuery({
    queryKey: ["device", id],
    queryFn: () => getDevice(id),
    refetchInterval: 30_000,
  });

  const systemInfoQuery = useQuery({
    queryKey: ["device-system-info", id],
    queryFn: () => getDeviceSystemInfo(id),
    enabled: activeTab === "system-info",
  });

  const softwareQuery = useQuery({
    queryKey: ["device-software", id],
    queryFn: () => getDeviceSoftware(id),
    enabled: activeTab === "software",
  });

  const processesQuery = useQuery({
    queryKey: ["device-processes", id],
    queryFn: () => getDeviceProcesses(id),
    enabled: activeTab === "processes",
  });

  const servicesQuery = useQuery({
    queryKey: ["device-services", id],
    queryFn: () => getDeviceServices(id),
    enabled: activeTab === "services",
  });

  const device = deviceQuery.data;

  if (deviceQuery.isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (deviceQuery.error || !device) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <AlertTriangle className="h-12 w-12 text-destructive mx-auto mb-2" />
          <p className="text-muted-foreground">Device not found</p>
          <Link href="/devices">
            <Button variant="outline" className="mt-4">
              <ArrowLeft className="mr-2 h-4 w-4" /> Back to Devices
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/devices">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
            <Monitor className="h-5 w-5 text-primary" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">{device.Hostname || device.DeviceName}</h1>
            <p className="text-sm text-muted-foreground">{device.IPAddress || "N/A"} · {device.OS || "N/A"}</p>
          </div>
        </div>
        <Badge variant={device.Status === "online" ? "success" : "secondary"} className="ml-auto capitalize">
          {device.Status === "online" ? "Online" : "Offline"}
        </Badge>
      </div>

      <div className="flex gap-1 border-b">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2 text-sm font-medium transition-colors border-b-2 -mb-px ${
              activeTab === tab.key
                ? "border-primary text-primary"
                : "border-transparent text-muted-foreground hover:text-foreground"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="min-h-[400px]">
        {activeTab === "overview" && <OverviewTab device={device} />}
        {activeTab === "live-view" && <LiveViewTab deviceId={id} />}
        {activeTab === "system-info" && (
          <SystemInfoTab query={systemInfoQuery} />
        )}
        {activeTab === "software" && (
          <SoftwareTab query={softwareQuery} />
        )}
        {activeTab === "processes" && (
          <ProcessesTab query={processesQuery} />
        )}
        {activeTab === "services" && (
          <ServicesTab deviceId={id} query={servicesQuery} />
        )}
      </div>
    </div>
  );
}
