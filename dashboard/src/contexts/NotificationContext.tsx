"use client";

import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";

export interface AppNotification {
  id: string;
  deviceId: string;
  message: string;
  timestamp: number;
}

interface NotificationContextType {
  notifications: AppNotification[];
  addNotification: (notification: Omit<AppNotification, "id" | "timestamp">) => void;
  removeNotification: (id: string) => void;
  clearAll: () => void;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const [notifications, setNotifications] = useState<AppNotification[]>([]);
  const router = useRouter();

  const addNotification = useCallback((notif: Omit<AppNotification, "id" | "timestamp">) => {
    const newNotif: AppNotification = {
      ...notif,
      id: Math.random().toString(36).substring(7),
      timestamp: Date.now(),
    };
    setNotifications((prev) => [newNotif, ...prev]);
  }, []);

  const removeNotification = useCallback((id: string) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setNotifications([]);
  }, []);

  // Auto-delete after 5 minutes
  useEffect(() => {
    const interval = setInterval(() => {
      const fiveMinsAgo = Date.now() - 5 * 60 * 1000;
      setNotifications((prev) => prev.filter((n) => n.timestamp > fiveMinsAgo));
    }, 10000); // Check every 10 seconds

    return () => clearInterval(interval);
  }, []);

  return (
    <NotificationContext.Provider value={{ notifications, addNotification, removeNotification, clearAll }}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useAppNotifications() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error("useAppNotifications must be used within a NotificationProvider");
  }
  return context;
}
