"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/hooks/use-auth";
import {
  LayoutDashboard,
  Monitor,
  Settings,
  User,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Activity,
  Grid,
} from "lucide-react";

export const navItems = [
  { href: "/", label: "Dashboard", icon: LayoutDashboard },
  { href: "/devices", label: "Devices", icon: Monitor },
  { href: "/monitoring", label: "Monitoring", icon: Activity },
  { href: "/multi-view", label: "Multi-View", icon: Grid },
  { href: "/settings", label: "Settings", icon: Settings },
  { href: "/profile", label: "Profile", icon: User },
];

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const pathname = usePathname();
  const { logout } = useAuth();

  return (
    <aside
      className={cn(
        "hidden md:flex flex-col border-r bg-sidebar text-sidebar-foreground transition-all duration-300",
        collapsed ? "w-16" : "w-60"
      )}
    >
      <div className="flex h-14 items-center gap-2 px-4">
        {!collapsed && (
          <span className="text-lg font-bold tracking-tight">
            Sentinel<span className="text-primary">Desk</span>
          </span>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setCollapsed(!collapsed)}
          className="ml-auto text-sidebar-foreground hover:bg-sidebar-muted"
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </Button>
      </div>

      <Separator className="bg-sidebar-border" />

      <nav className="flex-1 space-y-1 p-2">
        {navItems.map((item) => {
          const isActive = item.href === "/"
            ? pathname === "/"
            : pathname.startsWith(item.href);
          return (
            <Button
              key={item.href}
              variant="ghost"
              asChild
              className={cn(
                "w-full justify-start gap-3 text-sidebar-foreground hover:bg-sidebar-muted",
                isActive && "bg-sidebar-muted text-white",
                collapsed && "justify-center px-2"
              )}
            >
              <Link href={item.href}>
                <item.icon className="h-5 w-5 shrink-0" />
                {!collapsed && <span>{item.label}</span>}
              </Link>
            </Button>
          );
        })}
      </nav>

      <Separator className="bg-sidebar-border" />

      <div className="p-2">
        <Button
          variant="ghost"
          className={cn(
            "w-full justify-start gap-3 text-sidebar-foreground hover:bg-sidebar-muted hover:text-destructive",
            collapsed && "justify-center px-2"
          )}
          onClick={logout}
        >
          <LogOut className="h-5 w-5 shrink-0" />
          {!collapsed && <span>Logout</span>}
        </Button>
      </div>
    </aside>
  );
}
