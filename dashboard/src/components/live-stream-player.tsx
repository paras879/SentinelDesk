"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { AlertTriangle, Loader2, Monitor } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

interface Props {
  deviceId: string;
  deviceName?: string;
  onRemove?: () => void;
}

type ConnectionStatus = "disconnected" | "connecting" | "connected";

export function LiveStreamPlayer({ deviceId, deviceName, onRemove }: Props) {
  const [status, setStatus] = useState<ConnectionStatus>("connecting");
  const [displayFps, setDisplayFps] = useState(0);

  const wsRef = useRef<WebSocket | null>(null);
  const imgRef = useRef<HTMLImageElement>(null);

  const frameTimestamps = useRef<number[]>([]);
  const currentBlobUrl = useRef<string>("");
  const lastFpsUpdate = useRef(0);
  const reconnectAttempt = useRef(0);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const isUnmounted = useRef(false);

  const clearReconnect = useCallback(() => {
    if (reconnectTimer.current) {
      clearTimeout(reconnectTimer.current);
      reconnectTimer.current = undefined;
    }
  }, []);

  const connect = useCallback(() => {
    if (isUnmounted.current) return;
    
    const token = localStorage.getItem("sentineldesk_token");
    if (!token) {
      setStatus("disconnected");
      return;
    }

    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const wsBase = apiUrl.replace(/^http/, "ws");
    const url = `${wsBase}/ws/live/view/${deviceId}?token=${token}`;

    setStatus("connecting");
    const ws = new WebSocket(url);

    ws.onopen = () => {
      if (isUnmounted.current) {
        ws.close();
        return;
      }
      setStatus("connected");
      reconnectAttempt.current = 0;
    };

    ws.onmessage = (event) => {
      if (isUnmounted.current) return;
      
      if (event.data instanceof Blob) {
        const now = Date.now();

        frameTimestamps.current.push(now);
        if (frameTimestamps.current.length > 60) {
          frameTimestamps.current = frameTimestamps.current.slice(-60);
        }

        if (now - lastFpsUpdate.current >= 500) {
          const timestamps = frameTimestamps.current;
          if (timestamps.length > 1) {
            const elapsed = now - timestamps[0];
            const fps = Math.round(((timestamps.length - 1) / elapsed) * 1000);
            setDisplayFps(Math.min(fps, 60));
          }
          lastFpsUpdate.current = now;
        }

        if (currentBlobUrl.current) {
          URL.revokeObjectURL(currentBlobUrl.current);
        }
        const blob = new Blob([event.data], { type: "image/jpeg" });
        const blobUrl = URL.createObjectURL(blob);
        currentBlobUrl.current = blobUrl;
        if (imgRef.current) {
          imgRef.current.src = blobUrl;
        }
      }
    };

    ws.onclose = () => {
      if (isUnmounted.current) return;
      setStatus("disconnected");
      setDisplayFps(0);
      
      // Auto-reconnect logic
      const delay = Math.min(1000 * Math.pow(2, reconnectAttempt.current), 30000);
      reconnectAttempt.current++;
      reconnectTimer.current = setTimeout(connect, delay);
    };

    ws.onerror = () => {
      if (isUnmounted.current) return;
      setStatus("disconnected");
    };

    wsRef.current = ws;
  }, [deviceId]);

  useEffect(() => {
    isUnmounted.current = false;
    connect();

    return () => {
      isUnmounted.current = true;
      clearReconnect();
      if (currentBlobUrl.current) {
        URL.revokeObjectURL(currentBlobUrl.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect, clearReconnect]);

  return (
    <Card className="flex flex-col h-full overflow-hidden border-2 hover:border-primary/50 transition-colors">
      <div className="flex items-center justify-between p-2 px-3 bg-muted/50 border-b">
        <div className="flex items-center gap-2 truncate">
          <Monitor className="h-4 w-4 shrink-0 text-muted-foreground" />
          <span className="font-semibold text-sm truncate" title={deviceName || deviceId}>
            {deviceName || deviceId}
          </span>
        </div>
        <div className="flex items-center gap-3 shrink-0 text-xs">
          {status === "connected" ? (
            <>
              <span className="text-emerald-500 font-mono">{displayFps} FPS</span>
              <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
            </>
          ) : status === "connecting" ? (
            <>
              <span className="text-yellow-500">Connecting</span>
              <Loader2 className="h-3 w-3 animate-spin text-yellow-500" />
            </>
          ) : (
            <>
              <span className="text-destructive">Disconnected</span>
              <div className="h-2 w-2 rounded-full bg-destructive" />
            </>
          )}
          {onRemove && (
            <button
              onClick={onRemove}
              className="ml-1 text-muted-foreground hover:text-foreground p-1 rounded-md hover:bg-muted transition-colors"
              title="Remove from grid"
            >
              &times;
            </button>
          )}
        </div>
      </div>
      <CardContent className="p-0 flex-1 relative bg-black/95 min-h-[200px]">
        {status === "connected" ? (
          <img
            ref={imgRef}
            alt={`Live Stream ${deviceName || deviceId}`}
            className="w-full h-full object-contain"
            draggable={false}
          />
        ) : (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-muted-foreground">
            {status === "connecting" ? (
              <Loader2 className="h-8 w-8 animate-spin mb-2 opacity-50" />
            ) : (
              <AlertTriangle className="h-8 w-8 mb-2 opacity-50 text-destructive" />
            )}
            <span className="text-sm">
              {status === "connecting" ? "Establishing connection..." : "Connection lost"}
            </span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
