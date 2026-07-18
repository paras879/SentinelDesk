"use client";

import type { UseQueryResult } from "@tanstack/react-query";
import type { DeviceSystemInfo } from "@/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertTriangle } from "lucide-react";
import { formatBytes, formatUptime } from "@/lib/utils";

interface Props {
  query: UseQueryResult<DeviceSystemInfo, Error>;
}

export function SystemInfoTab({ query }: Props) {
  const { data, isLoading, error } = query;

  if (isLoading) {
    return (
      <div className="grid gap-6 md:grid-cols-2">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader><Skeleton className="h-5 w-32" /></CardHeader>
            <CardContent className="space-y-3">
              {[...Array(4)].map((_, j) => (
                <Skeleton key={j} className="h-4 w-full" />
              ))}
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <AlertTriangle className="h-12 w-12 text-destructive mx-auto mb-2" />
          <p className="text-muted-foreground">System info not available</p>
        </div>
      </div>
    );
  }

  return (
    <div className="grid gap-6 md:grid-cols-2">
      <Card>
        <CardHeader><CardTitle className="text-lg">CPU</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <SysRow label="Model" value={data.cpu_name || "—"} />
          <SysRow label="Usage" value={data.cpu_usage != null ? `${data.cpu_usage.toFixed(1)}%` : "—"} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle className="text-lg">Memory</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <SysRow label="Total" value={data.total_ram != null ? formatBytes(data.total_ram) : "—"} />
          <SysRow label="Used" value={data.used_ram != null ? formatBytes(data.used_ram) : "—"} />
          <SysRow label="Free" value={data.free_ram != null ? formatBytes(data.free_ram) : "—"} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle className="text-lg">Disk</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <SysRow label="Total" value={data.total_disk != null ? formatBytes(data.total_disk) : "—"} />
          <SysRow label="Used" value={data.used_disk != null ? formatBytes(data.used_disk) : "—"} />
          <SysRow label="Free" value={data.free_disk != null ? formatBytes(data.free_disk) : "—"} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle className="text-lg">Network</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <SysRow label="Local IP" value={data.local_ip || "—"} />
          <SysRow label="Public IP" value={data.public_ip || "—"} />
          <SysRow label="MAC Address" value={data.mac_address || "—"} />
        </CardContent>
      </Card>

      <Card className="md:col-span-2">
        <CardHeader><CardTitle className="text-lg">System Details</CardTitle></CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-2">
          <SysRow label="OS Version" value={data.os_version || "—"} />
          <SysRow label="Agent Version" value={data.agent_version || "—"} />
          <SysRow label="Uptime" value={data.uptime ? formatUptime(data.uptime) : "—"} />
          <SysRow label="Last Boot" value={data.last_boot ? new Date(data.last_boot).toLocaleString() : "—"} />
        </CardContent>
      </Card>
    </div>
  );
}

function SysRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className="text-sm font-medium">{value}</span>
    </div>
  );
}
