"use client";

import { useState, useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { getDevices, deleteDevice, updateDeviceLocation } from "@/services/devices";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
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
  Trash2,
  User,
  Home,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type SiteFilter = "all" | "WFH" | "Office";

export default function DevicesPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [siteFilter, setSiteFilter] = useState<SiteFilter>("all");
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [updatingId, setUpdatingId] = useState<string | null>(null);

  const locationType = siteFilter === "all" ? undefined : siteFilter;

  const { data, isLoading, error } = useQuery({
    queryKey: ["devices", locationType],
    queryFn: () => getDevices(locationType),
    refetchInterval: 5_000,
  });

  const filtered = data?.devices?.filter(
    (d) =>
      d.Hostname?.toLowerCase().includes(search.toLowerCase()) ||
      d.ID?.toLowerCase().includes(search.toLowerCase()) ||
      d.IPAddress?.toLowerCase().includes(search.toLowerCase()) ||
      d.ConnectedSubnet?.toLowerCase().includes(search.toLowerCase())
  );

  const handleDelete = async (id: string, name: string) => {
    if (!window.confirm(`Are you sure you want to delete the device ${name}? This action cannot be undone.`)) {
      return;
    }

    try {
      setDeletingId(id);
      await deleteDevice(id);
      toast.success("Device deleted successfully");
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    } catch (error: any) {
      toast.error(error?.response?.data?.error || "Failed to delete device");
    } finally {
      setDeletingId(null);
    }
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
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Devices</h1>
          <p className="text-sm text-muted-foreground">
            {isLoading ? "Loading..." : `${data?.count || 0} devices`}
                      </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            size="sm"
            variant={siteFilter === "WFH" ? "default" : "outline"}
            onClick={() => setSiteFilter("WFH")}
          >
            <Home className="mr-1 h-4 w-4" />
            WFH
          </Button>
          <Button
            size="sm"
            variant={siteFilter === "Office" ? "default" : "outline"}
            onClick={() => setSiteFilter("Office")}
          >
            <Building2 className="mr-1 h-4 w-4" />
            Office
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
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <CardTitle className="text-lg">
              {siteFilter === "all" ? "All Devices" : `${siteFilter} Devices`}
            </CardTitle>
            <div className="relative w-full sm:w-72">
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
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="whitespace-nowrap">Status</TableHead>
                    <TableHead className="whitespace-nowrap">Device Name</TableHead>
                    <TableHead className="whitespace-nowrap">Device ID</TableHead>
                    <TableHead className="whitespace-nowrap">Current IP</TableHead>
                    <TableHead className="whitespace-nowrap">Site Group</TableHead>
                    <TableHead className="whitespace-nowrap">Location</TableHead>
                    <TableHead className="whitespace-nowrap">User</TableHead>
                    <TableHead className="whitespace-nowrap">Last Seen</TableHead>
                    <TableHead className="text-right whitespace-nowrap">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filtered?.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                        {search
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
                        <TableCell className="font-medium whitespace-nowrap">{device.DeviceName || device.Hostname || "—"}</TableCell>
                        <TableCell className="font-mono text-xs text-muted-foreground whitespace-nowrap">
                          {device.ID.slice(0, 16)}...
                        </TableCell>
                        <TableCell className="whitespace-nowrap">{device.IPAddress || "—"}</TableCell>
                        <TableCell className="whitespace-nowrap">
                          {device.NetworkGroupID ? (
                            <Badge variant="outline" className="text-xs font-mono">
                              <Layers className="mr-1 h-3 w-3" />
                              {device.NetworkGroupID.slice(0, 8)}...
                            </Badge>
                          ) : (
                            <span className="text-xs text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell className="whitespace-nowrap">
                          <Select
                            disabled={updatingId === device.ID}
                            value={device.LocationType || "Unassigned"}
                            onValueChange={async (val) => {
                              try {
                                setUpdatingId(device.ID);
                                await updateDeviceLocation(device.ID, val);
                                queryClient.invalidateQueries({ queryKey: ["devices"] });
                                toast.success("Location updated");
                              } catch {
                                toast.error("Failed to update location");
                              } finally {
                                setUpdatingId(null);
                              }
                            }}
                          >
                            <SelectTrigger className="h-8 w-[130px] text-xs">
                              <SelectValue placeholder="Location" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="Unassigned">Unassigned</SelectItem>
                              <SelectItem value="WFH">WFH</SelectItem>
                              <SelectItem value="Office">Office</SelectItem>
                              <SelectItem value="Admin">Admin</SelectItem>
                            </SelectContent>
                          </Select>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
                          {device.Username ? (
                            <div className="flex items-center gap-1.5" title={device.Username}>
                              <User className="h-3.5 w-3.5 shrink-0" />
                              <span className="truncate max-w-[120px]">{device.Username.split('\\').pop()}</span>
                            </div>
                          ) : (
                            "—"
                          )}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
                          {device.LastSeen
                            ? new Date(device.LastSeen).toLocaleString()
                            : "Never"}
                        </TableCell>
                        <TableCell className="text-right whitespace-nowrap">
                          <div className="flex items-center justify-end gap-2">
                            <Button
                              variant="link"
                              size="sm"
                              onClick={() => router.push(`/devices/${device.ID}`)}
                              className="h-auto p-0"
                            >
                              <Eye className="h-3.5 w-3.5 mr-1" />
                              View
                            </Button>
                            <Button
                              variant="link"
                              size="sm"
                              onClick={() => handleDelete(device.ID, device.DeviceName || device.Hostname || device.ID)}
                              className="h-auto p-0 text-destructive hover:text-destructive/80"
                              disabled={deletingId === device.ID}
                            >
                              <Trash2 className="h-3.5 w-3.5 mr-1" />
                              Delete
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
