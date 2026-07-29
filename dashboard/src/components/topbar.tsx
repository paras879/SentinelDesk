"use client";

import { Bell, Search, User, LogOut, Settings, Menu, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme-toggle";
import { Sheet, SheetContent, SheetTrigger, SheetTitle, SheetHeader } from "@/components/ui/sheet";
import { navItems } from "@/components/sidebar";
import { usePathname, useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { useAuth } from "@/hooks/use-auth";
import { useAppNotifications } from "@/contexts/NotificationContext";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import Link from "next/link";

export function Topbar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const router = useRouter();
  const { notifications, removeNotification } = useAppNotifications();

  return (
    <header className="flex h-14 items-center gap-4 border-b bg-background px-4 md:px-6">
      <div className="flex items-center md:hidden">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" className="-ml-2 mr-2 text-muted-foreground">
              <Menu className="h-5 w-5" />
              <span className="sr-only">Toggle Menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-72 p-0 flex flex-col">
            <SheetHeader className="h-14 border-b px-6 flex items-center justify-start flex-row mt-0">
              <SheetTitle className="text-lg font-bold tracking-tight text-left mt-0">
                Sentinel<span className="text-primary">Desk</span>
              </SheetTitle>
            </SheetHeader>
            <nav className="flex-1 space-y-1 p-4">
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
                      "w-full justify-start gap-3",
                      isActive && "bg-muted text-foreground"
                    )}
                  >
                    <Link href={item.href}>
                      <item.icon className="h-5 w-5 shrink-0" />
                      <span>{item.label}</span>
                    </Link>
                  </Button>
                );
              })}
            </nav>
          </SheetContent>
        </Sheet>
      </div>

      <div className="relative flex-1 max-w-md hidden sm:block">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input placeholder="Search devices..." className="pl-9" />
      </div>

      <div className="flex items-center gap-4 ml-auto">
        <ThemeToggle />
        
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="ghost" size="icon" className="relative">
              <Bell className="h-5 w-5" />
              {notifications.length > 0 && (
                <span className="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] font-medium text-primary-foreground">
                  {notifications.length}
                </span>
              )}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-80" align="end">
            <div className="grid gap-4">
              <div className="space-y-2">
                <h4 className="font-medium leading-none">Notifications</h4>
                <p className="text-sm text-muted-foreground">
                  {notifications.length === 0 
                    ? "You have no unread messages." 
                    : `You have ${notifications.length} unread message${notifications.length === 1 ? '' : 's'}.`}
                </p>
              </div>
              <div className="grid gap-2 max-h-[300px] overflow-y-auto pr-1">
                {notifications.length === 0 ? (
                  <div className="text-sm text-muted-foreground py-2 text-center border-t pt-4">No new notifications</div>
                ) : (
                  notifications.map((notif) => (
                    <div key={notif.id} className="group relative flex items-center justify-between border-b pb-2 pt-1 last:border-0 last:pb-0">
                      <div 
                        className="text-sm pr-6 cursor-pointer hover:text-primary transition-colors flex-1"
                        onClick={() => router.push(`/devices/${notif.deviceId}`)}
                      >
                        {notif.message}
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={(e) => {
                          e.stopPropagation();
                          removeNotification(notif.id);
                        }}
                      >
                        <X className="h-4 w-4 text-muted-foreground hover:text-foreground" />
                      </Button>
                    </div>
                  ))
                )}
              </div>
            </div>
          </PopoverContent>
        </Popover>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-8 w-8 rounded-full">
              <Avatar className="h-8 w-8">
                <AvatarFallback>{user?.name?.charAt(0)?.toUpperCase() || "A"}</AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56" align="end" forceMount>
            <DropdownMenuLabel className="font-normal">
              <div className="flex flex-col space-y-1">
                <p className="text-sm font-medium leading-none">{user?.name || "Admin"}</p>
                <p className="text-xs leading-none text-muted-foreground">
                  {user?.email || "admin@sentineldesk.com"}
                </p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link href="/profile" className="cursor-pointer w-full flex items-center">
                <User className="mr-2 h-4 w-4" />
                <span>Profile</span>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link href="/settings" className="cursor-pointer w-full flex items-center">
                <Settings className="mr-2 h-4 w-4" />
                <span>Settings</span>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer" onClick={() => logout && logout()}>
              <LogOut className="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
