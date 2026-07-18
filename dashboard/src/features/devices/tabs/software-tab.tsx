"use client";

import { useState } from "react";
import type { UseQueryResult } from "@tanstack/react-query";
import type { SoftwareItem } from "@/types";
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
import { AlertTriangle, Search, Package } from "lucide-react";
import { formatBytes } from "@/lib/utils";

interface Props {
  query: UseQueryResult<{ count: number; software: SoftwareItem[] }, Error>;
}

export function SoftwareTab({ query }: Props) {
  const [search, setSearch] = useState("");
  const { data, isLoading, error } = query;

  const filtered = data?.software?.filter(
    (s) =>
      s.name?.toLowerCase().includes(search.toLowerCase()) ||
      s.publisher?.toLowerCase().includes(search.toLowerCase()) ||
      s.version?.toLowerCase().includes(search.toLowerCase())
  );

  if (isLoading) {
    return <div className="space-y-3">{[...Array(5)].map((_, i) => <Skeleton key={i} className="h-10 w-full" />)}</div>;
  }

  if (error || !data) {
    return (
      <div className="flex items-center justify-center h-64">
        <AlertTriangle className="h-8 w-8 text-destructive mr-2" />
        <p className="text-muted-foreground">Software data not available</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">{data.count || data.software?.length || 0} applications</p>
        <div className="relative w-64">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="Search software..." className="pl-9" value={search} onChange={(e) => setSearch(e.target.value)} />
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead><Package className="h-4 w-4" /></TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Version</TableHead>
              <TableHead>Publisher</TableHead>
              <TableHead>Install Date</TableHead>
              <TableHead>Size</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filtered?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {search ? "No matches" : "No software found"}
                </TableCell>
              </TableRow>
            ) : (
              filtered?.map((item, i) => (
                <TableRow key={item.ID || i}>
                  <TableCell><Package className="h-4 w-4 text-muted-foreground" /></TableCell>
                  <TableCell className="font-medium">{item.name || "—"}</TableCell>
                  <TableCell className="text-muted-foreground">{item.version || "—"}</TableCell>
                  <TableCell className="text-muted-foreground">{item.publisher || "—"}</TableCell>
                  <TableCell className="text-muted-foreground text-sm">
                    {item.install_date ? new Date(item.install_date).toLocaleDateString() : "—"}
                  </TableCell>
                  <TableCell className="text-muted-foreground text-sm">
                    {item.estimated_size != null ? formatBytes(item.estimated_size) : "—"}
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
