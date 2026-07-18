"use client";

import type { Device } from "@/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";

interface Props {
  device: Device;
}

export function OverviewTab({ device }: Props) {
  const adapters = parseNetworkAdapters(device.NetworkAdapters);

  return (
    <div className="grid gap-6 md:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Device Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Row label="Device ID" value={device.ID} mono />
          <Row label="Device Name" value={device.DeviceName || "—"} />
          <Row label="Hostname" value={device.Hostname || "—"} />
          <Row label="Logged-in User" value={device.Username || "—"} />
          <Row label="Current IP" value={device.IPAddress || "—"} />
          <Row label="MAC Address" value={device.MACAddress || "—"} />
          <Row label="Connected Network" value={device.ConnectedSubnet || "—"} />
          <Row label="Status">
            <Badge variant={device.Status === "online" ? "success" : "secondary"}>
              {device.Status === "online" ? "Online" : "Offline"}
            </Badge>
          </Row>
          <Row label="Last Seen" value={device.LastSeen ? formatDate(device.LastSeen) : "Never"} />
          <Row label="Registered" value={device.CreatedAt ? formatDate(device.CreatedAt) : "—"} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Operating System</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Row label="OS" value={device.OS || "—"} />
          {device.OSVersion && <Row label="Version" value={device.OSVersion} />}
        </CardContent>
      </Card>

      {adapters.length > 0 && (
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle className="text-lg">Network Adapters</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {adapters.map((adapter, i) => (
                <div key={i} className="border rounded-lg p-3 space-y-1 text-sm">
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{adapter.name}</span>
                    {adapter.is_primary && (
                      <Badge variant="outline" className="text-[10px]">Primary</Badge>
                    )}
                  </div>
                  <div className="text-muted-foreground">
                    <span>MAC: {adapter.mac_address}</span>
                    <span className="ml-4">IP: {adapter.ip_addresses}</span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

interface NetworkAdapter {
  name: string;
  mac_address: string;
  ip_addresses: string;
  subnet_mask: string;
  is_primary: boolean;
}

function parseNetworkAdapters(json: string | undefined): NetworkAdapter[] {
  if (!json) return [];
  try {
    return JSON.parse(json);
  } catch {
    return [];
  }
}

function Row({ label, value, children, mono }: { label: string; value?: string; children?: React.ReactNode; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-sm text-muted-foreground">{label}</span>
      {children || (
        <span className={`text-sm font-medium ${mono ? "font-mono text-xs" : ""}`}>
          {value}
        </span>
      )}
    </div>
  );
}
