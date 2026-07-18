"use client";

import { useState } from "react";
import type { UseQueryResult } from "@tanstack/react-query";
import type { ProcessItem } from "@/types";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AlertTriangle, Search, Cpu } from "lucide-react";
import { formatBytes } from "@/lib/utils";

interface Props {
  query: UseQueryResult<{ count: number; processes: ProcessItem[] }, Error>;
}

export function ProcessesTab({ query }: Props) {
  const [search, setSearch] = useState("");
  const { data, isLoading, error } = query;

  const filtered = data?.processes?.filter(
    (p) =>
      p.name?.toLowerCase().includes(search.toLowerCase()) ||
      String(p.pid).includes(search)
  );

  if (isLoading) {
    return <div className="space-y-3">{[...Array(5)].map((_, i) => <Skeleton key={i} className="h-10 w-full" />)}</div>;
  }

  if (error || !data) {
    return (
      <div className="flex items-center justify-center h-64">
        <AlertTriangle className="h-8 w-8 text-destructive mr-2" />
        <p className="text-muted-foreground">Process data not available</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">{data.count || data.processes?.length || 0} processes</p>
        <div className="relative w-64">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="Search processes..." className="pl-9" value={search} onChange={(e) => setSearch(e.target.value)} />
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>PID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>CPU %</TableHead>
              <TableHead>Memory</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>User</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filtered?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {search ? "No matches" : "No processes found"}
                </TableCell>
              </TableRow>
            ) : (
              filtered?.map((proc, i) => (
                <TableRow key={proc.pid || i}>
                  <TableCell className="font-mono text-xs text-muted-foreground">{proc.pid}</TableCell>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      <Cpu className="h-3.5 w-3.5 text-muted-foreground" />
                      {proc.name || "—"}
                    </div>
                  </TableCell>
                  <TableCell>{proc.cpu_usage != null ? `${proc.cpu_usage.toFixed(1)}%` : "—"}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {proc.memory_bytes != null ? formatBytes(proc.memory_bytes) : "—"}
                  </TableCell>
                  <TableCell className="text-sm">{proc.status || "—"}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">{proc.username || "—"}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
