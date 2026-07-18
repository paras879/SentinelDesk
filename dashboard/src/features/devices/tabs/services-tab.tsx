"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { UseQueryResult } from "@tanstack/react-query";
import type { WindowsServiceItem } from "@/types";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AlertTriangle, Search, Play, Square, RefreshCw, Loader2 } from "lucide-react";
import { startService, stopService, restartService } from "@/services/devices";
import { toast } from "sonner";

interface Props {
  deviceId: string;
  query: UseQueryResult<{ count: number; services: WindowsServiceItem[] }, Error>;
}

export function ServicesTab({ deviceId, query }: Props) {
  const [search, setSearch] = useState("");
  const { data, isLoading, error } = query;
  const queryClient = useQueryClient();

  const invalidate = () => queryClient.invalidateQueries({ queryKey: ["device-services", deviceId] });

  const startMut = useMutation({
    mutationFn: (name: string) => startService(deviceId, name),
    onSuccess: () => { toast.success("Service started"); invalidate(); },
    onError: (e: any) => toast.error(e?.response?.data?.error || "Failed to start"),
  });

  const stopMut = useMutation({
    mutationFn: (name: string) => stopService(deviceId, name),
    onSuccess: () => { toast.success("Service stopped"); invalidate(); },
    onError: (e: any) => toast.error(e?.response?.data?.error || "Failed to stop"),
  });

  const restartMut = useMutation({
    mutationFn: (name: string) => restartService(deviceId, name),
    onSuccess: () => { toast.success("Service restarted"); invalidate(); },
    onError: (e: any) => toast.error(e?.response?.data?.error || "Failed to restart"),
  });

  const filtered = data?.services?.filter(
    (s) =>
      s.service_name?.toLowerCase().includes(search.toLowerCase()) ||
      s.display_name?.toLowerCase().includes(search.toLowerCase())
  );

  if (isLoading) {
    return <div className="space-y-3">{[...Array(5)].map((_, i) => <Skeleton key={i} className="h-10 w-full" />)}</div>;
  }

  if (error || !data) {
    return (
      <div className="flex items-center justify-center h-64">
        <AlertTriangle className="h-8 w-8 text-destructive mr-2" />
        <p className="text-muted-foreground">Service data not available</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">{data.count || data.services?.length || 0} services</p>
        <div className="relative w-64">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="Search services..." className="pl-9" value={search} onChange={(e) => setSearch(e.target.value)} />
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Display Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Start Type</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filtered?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                  {search ? "No matches" : "No services found"}
                </TableCell>
              </TableRow>
            ) : (
              filtered?.map((svc, i) => (
                <TableRow key={svc.service_name || i}>
                  <TableCell className="font-mono text-xs">{svc.service_name}</TableCell>
                  <TableCell className="font-medium">{svc.display_name || svc.service_name}</TableCell>
                  <TableCell>
                    <Badge variant={svc.status === "Running" || svc.status === "running" ? "success" : "secondary"}>
                      {svc.status || "—"}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">{svc.start_type || "—"}</TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end gap-1">
                      {(svc.status === "Stopped" || svc.status === "stopped") && (
                        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => startMut.mutate(svc.service_name)} disabled={startMut.isPending}>
                          {startMut.isPending ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Play className="h-3.5 w-3.5 text-emerald-500" />}
                        </Button>
                      )}
                      {(svc.status === "Running" || svc.status === "running") && (
                        <>
                          <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => stopMut.mutate(svc.service_name)} disabled={stopMut.isPending}>
                            {stopMut.isPending ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Square className="h-3.5 w-3.5 text-destructive" />}
                          </Button>
                          <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => restartMut.mutate(svc.service_name)} disabled={restartMut.isPending}>
                            {restartMut.isPending ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <RefreshCw className="h-3.5 w-3.5" />}
                          </Button>
                        </>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
