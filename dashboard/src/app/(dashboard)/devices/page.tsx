"use client";

import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { getDevices } from "@/services/devices";
import { getMyIP } from "@/services/dashboard";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Search,
  Monitor,
  AlertTriangle,
  Wifi,
  WifiOff,
  Eye,
  Globe,
  Building2,
  Layers,
} from "lucide-react";

type SiteFilter = "all" | "current";

export default function DevicesPage() {
  const router = useRouter();
  const [search, setSearch] = useState("");
  const [siteFilter, setSiteFilter] = useState<SiteFilter>("current");
  const [mySiteGroup, setMySiteGroup] = useState<string>("");

  useEffect(() => {
    (async () => {
      try {
        const res = await getMyIP();
        if (res.network_group_id) {
          setMySiteGroup(res.network_group_id);
        } else {
          setMySiteGroup("");
          setSiteFilter("all");
        }
      } catch {
        setMySiteGroup("");
        setSiteFilter("all");
      }
    })();
  }, []);

  const networkGroupID = siteFilter === "current" ? mySiteGroup : undefined;

  const { data, isLoading, error } = useQuery({
    queryKey: ["devices", networkGroupID],
    queryFn: () => getDevices(networkGroupID),
    refetchInterval: 15_000,
    enabled: siteFilter === "all" || (siteFilter === "current" && !!mySiteGroup) || siteFilter === "current",
  });

  const filtered = data?.devices?.filter(
    (d) =>
      d.Hostname?.toLowerCase().includes(search.toLowerCase()) ||
      d.ID?.toLowerCase().includes(search.toLowerCase()) ||
      d.IPAddress?.toLowerCase().includes(search.toLowerCase()) ||
      d.ConnectedSubnet?.toLowerCase().includes(search.toLowerCase())
  );

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
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Devices</h1>
          <p className="text-sm text-muted-foreground">
            {isLoading ? "Loading..." : `${data?.count || 0} devices`}
            {siteFilter === "current" && mySiteGroup && (
              <span className="ml-2 font-mono text-xs">(Site Group: {mySiteGroup.slice(0, 8)}...)</span>
            )}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            size="sm"
            variant={siteFilter === "current" ? "default" : "outline"}
            onClick={() => setSiteFilter("current")}
            disabled={!mySiteGroup}
          >
            <Building2 className="mr-1 h-4 w-4" />
            Current Site
          </Button>
          <Button
            size="sm"
            variant={siteFilter === "all" ? "default" : "outline"}
            onClick={() => setSiteFilter("all")}
          >
            <Globe className="mr-1 h-4 w-4" />
            All Devices
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg">
              {siteFilter === "current" ? "Current Site Devices" : "All Devices"}
            </CardTitle>
            <div className="relative w-72">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search by name, IP, network..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Status</TableHead>
                  <TableHead>Device Name</TableHead>
                  <TableHead>Device ID</TableHead>
                  <TableHead>Current IP</TableHead>
                  <TableHead>Site Group</TableHead>
                  <TableHead>User</TableHead>
                  <TableHead>Last Seen</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered?.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                      {siteFilter === "current" && !mySiteGroup
                        ? "Unable to detect your site. Switch to All Devices."
                        : search
                          ? "No devices match your search"
                          : "No devices registered yet"}
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered?.map((device) => (
                    <TableRow key={device.ID}>
                      <TableCell>
                        {device.Status === "online" ? (
                          <Wifi className="h-4 w-4 text-emerald-500" />
                        ) : (
                          <WifiOff className="h-4 w-4 text-muted-foreground" />
                        )}
                      </TableCell>
                      <TableCell className="font-medium">{device.DeviceName || device.Hostname || "—"}</TableCell>
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {device.ID.slice(0, 16)}...
                      </TableCell>
                      <TableCell>{device.IPAddress || "—"}</TableCell>
                      <TableCell>
                        {device.NetworkGroupID ? (
                          <Badge variant="outline" className="text-xs font-mono">
                            <Layers className="mr-1 h-3 w-3" />
                            {device.NetworkGroupID.slice(0, 8)}...
                          </Badge>
                        ) : (
                          <span className="text-xs text-muted-foreground">—</span>
                        )}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {device.Username || "—"}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {device.LastSeen
                          ? new Date(device.LastSeen).toLocaleString()
                          : "Never"}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="link"
                          size="sm"
                          onClick={() => router.push(`/devices/${device.ID}`)}
                          className="h-auto p-0"
                        >
                          <Eye className="h-3.5 w-3.5 mr-1" />
                          View
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
