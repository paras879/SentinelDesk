"use client";

import { useQuery } from "@tanstack/react-query";
import { getDevices } from "@/services/devices";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Activity, AlertTriangle, Monitor, Server, WifiOff, User } from "lucide-react";

export default function MonitoringPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["devices", "monitoring"],
    queryFn: () => getDevices(),
    refetchInterval: 5_000,
  });

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <AlertTriangle className="h-12 w-12 text-destructive mx-auto mb-2" />
          <p className="text-muted-foreground">Failed to load monitoring data</p>
        </div>
      </div>
    );
  }

  const devices = data?.devices || [];
  const onlineDevices = devices.filter((d: any) => d.Status === "online");
  const offlineDevices = devices.filter((d: any) => d.Status !== "online");
  
  // Mock warnings for NOC feel
  const warnings = offlineDevices.length + (devices.length > 5 ? 2 : 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold tracking-tight">NOC Board</h1>
        <p className="text-sm text-muted-foreground">Real-time status of all endpoints and agents.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Devices</CardTitle>
            <Server className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold">{devices.length}</div>
            )}
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Online</CardTitle>
            <Activity className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold text-emerald-500">{onlineDevices.length}</div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Offline</CardTitle>
            <WifiOff className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold">{offlineDevices.length}</div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Alerts</CardTitle>
            <AlertTriangle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold text-destructive">{warnings}</div>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {isLoading ? (
          Array.from({ length: 6 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="p-6 flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <Skeleton className="h-12 w-12 rounded-full" />
                  <div className="space-y-2">
                    <Skeleton className="h-4 w-[100px]" />
                    <Skeleton className="h-4 w-[80px]" />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        ) : devices.length === 0 ? (
          <div className="col-span-full py-12 text-center text-muted-foreground">
            No devices found on the network.
          </div>
        ) : (
          devices.map((device: any) => {
            const isOnline = device.Status === "online";
            return (
              <Card key={device.ID} className={`border-l-4 ${isOnline ? 'border-l-emerald-500' : 'border-l-destructive'}`}>
                <CardContent className="p-4">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center space-x-3">
                      <div className={`p-2 rounded-md ${isOnline ? 'bg-emerald-500/10 text-emerald-500' : 'bg-destructive/10 text-destructive'}`}>
                        <Monitor className="h-6 w-6" />
                      </div>
                      <div>
                        <h3 className="font-semibold truncate w-32 md:w-40" title={device.DeviceName || device.Hostname}>
                          {device.DeviceName || device.Hostname || "Unknown Device"}
                        </h3>
                        <div className="flex items-center gap-1.5 mt-2 mb-1 w-full" title={device.Username}>
                          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-secondary/60 text-secondary-foreground text-sm font-medium border shadow-sm max-w-[140px] md:max-w-[160px]">
                            <User className="h-4 w-4 shrink-0 opacity-70" />
                            <span className="truncate">{device.Username ? device.Username.split('\\').pop() : "Unknown User"}</span>
                          </div>
                        </div>
                        <p className="text-xs text-muted-foreground font-mono">
                          {device.IPAddress || "No IP"}
                        </p>
                      </div>
                    </div>
                    {!isOnline && (
                      <AlertTriangle className="h-5 w-5 text-destructive" />
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })
        )}
      </div>
    </div>
  );
}
