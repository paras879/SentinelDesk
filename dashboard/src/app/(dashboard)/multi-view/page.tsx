"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getDevices } from "@/services/devices";
import { LiveStreamPlayer } from "@/components/live-stream-player";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Grid, Monitor, AlertTriangle, MonitorPlay } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

const MAX_DEVICES = 6;

export default function MultiViewPage() {
  const [selectedDeviceIds, setSelectedDeviceIds] = useState<string[]>([]);
  
  const { data, isLoading, error } = useQuery({
    queryKey: ["devices", "multi-view"],
    queryFn: () => getDevices(),
    refetchInterval: 10_000,
  });

  const onlineDevices = data?.devices?.filter((d: any) => d.Status === "online") || [];

  const handleToggleDevice = (deviceId: string) => {
    setSelectedDeviceIds((prev) => {
      if (prev.includes(deviceId)) {
        return prev.filter((id) => id !== deviceId);
      }
      if (prev.length >= MAX_DEVICES) {
        return prev;
      }
      return [...prev, deviceId];
    });
  };

  const handleRemoveDevice = (deviceId: string) => {
    setSelectedDeviceIds((prev) => prev.filter((id) => id !== deviceId));
  };

  const handleClearAll = () => {
    setSelectedDeviceIds([]);
  };

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <AlertTriangle className="h-12 w-12 text-destructive mx-auto mb-2" />
          <p className="text-muted-foreground">Failed to load devices</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 h-full flex flex-col">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Grid className="h-6 w-6" />
            Multi-Live View
          </h1>
          <p className="text-sm text-muted-foreground">
            Monitor up to {MAX_DEVICES} devices simultaneously. (Selected: {selectedDeviceIds.length}/{MAX_DEVICES})
          </p>
        </div>
        {selectedDeviceIds.length > 0 && (
          <Button variant="outline" size="sm" onClick={handleClearAll}>
            Clear All
          </Button>
        )}
      </div>

      <div className="flex flex-col xl:flex-row gap-6 flex-1 min-h-0">
        {/* Device Selection Sidebar */}
        <Card className="xl:w-80 flex flex-col shrink-0">
          <CardHeader className="py-4 border-b">
            <CardTitle className="text-base font-semibold flex items-center gap-2">
              <MonitorPlay className="h-4 w-4" />
              Available Devices
            </CardTitle>
          </CardHeader>
          <CardContent className="p-0 flex-1 overflow-y-auto">
            {isLoading ? (
              <div className="p-4 space-y-3">
                {[...Array(5)].map((_, i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            ) : onlineDevices.length === 0 ? (
              <div className="p-6 text-center text-muted-foreground text-sm">
                No online devices available.
              </div>
            ) : (
              <div className="divide-y">
                {onlineDevices.map((device: any) => {
                  const isSelected = selectedDeviceIds.includes(device.ID);
                  const isDisabled = !isSelected && selectedDeviceIds.length >= MAX_DEVICES;

                  return (
                    <div
                      key={device.ID}
                      className={`flex items-center space-x-3 p-3 hover:bg-muted/50 transition-colors ${
                        isDisabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"
                      }`}
                      onClick={() => !isDisabled && handleToggleDevice(device.ID)}
                    >
                      <Checkbox
                        checked={isSelected}
                        disabled={isDisabled}
                        onCheckedChange={() => handleToggleDevice(device.ID)}
                      />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium leading-none truncate">
                          {device.DeviceName || device.Hostname || "Unknown"}
                        </p>
                        <p className="text-xs text-muted-foreground truncate mt-1">
                          {device.Username ? device.Username.split('\\').pop() : "No User"}
                        </p>
                      </div>
                      <div className="h-2 w-2 rounded-full bg-emerald-500 shrink-0" />
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Video Grid */}
        <div className="flex-1 min-h-[500px] xl:min-h-0 bg-muted/30 border rounded-xl p-4 overflow-y-auto">
          {selectedDeviceIds.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-muted-foreground space-y-4">
              <Grid className="h-16 w-16 opacity-20" />
              <p>Select devices from the list to start monitoring.</p>
            </div>
          ) : (
            <div
              className={`grid gap-4 h-full ${
                selectedDeviceIds.length === 1
                  ? "grid-cols-1 grid-rows-1"
                  : selectedDeviceIds.length <= 2
                  ? "grid-cols-1 md:grid-cols-2 grid-rows-1"
                  : selectedDeviceIds.length <= 4
                  ? "grid-cols-1 md:grid-cols-2 grid-rows-2"
                  : "grid-cols-1 md:grid-cols-2 lg:grid-cols-3 grid-rows-2"
              }`}
            >
              {selectedDeviceIds.map((id) => {
                const device = onlineDevices.find((d: any) => d.ID === id);
                return (
                  <div key={id} className="min-h-[250px]">
                    <LiveStreamPlayer
                      deviceId={id}
                      deviceName={device?.DeviceName || device?.Hostname}
                      onRemove={() => handleRemoveDevice(id)}
                    />
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
