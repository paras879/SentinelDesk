"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Play, Square, RefreshCw, Monitor, MousePointer } from "lucide-react";
import { useRemoteControl } from "@/hooks/use-remote-control";

interface Props {
  deviceId: string;
}

type ConnectionStatus = "disconnected" | "connecting" | "connected";

export function LiveViewTab({ deviceId }: Props) {
  const [status, setStatus] = useState<ConnectionStatus>("disconnected");
  const [displayFps, setDisplayFps] = useState(0);
  const [displayResolution, setDisplayResolution] = useState("");
  const [lastFrameTime, setLastFrameTime] = useState("");

  const wsRef = useRef<WebSocket | null>(null);
  const imgRef = useRef<HTMLImageElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const frameTimestamps = useRef<number[]>([]);
  const currentBlobUrl = useRef<string>("");
  const frameCount = useRef(0);
  const lastFpsUpdate = useRef(0);
  const lastSeq = useRef(0);

  const isManualDisconnect = useRef(false);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const reconnectAttempt = useRef(0);

  const {
    enabled: remoteEnabled,
    enable: enableRemote,
    disable: disableRemote,
    screenInfo,
    cursorX,
    cursorY,
    packetsSent,
    sendMouseMove,
    sendMouseDown,
    sendMouseUp,
    sendDoubleClick,
    sendMouseWheel,
    handleWSMessage,
    handleKeyEvent: hookHandleKeyEvent,
  } = useRemoteControl(wsRef, containerRef);

  const clearReconnect = useCallback(() => {
    if (reconnectTimer.current) {
      clearTimeout(reconnectTimer.current);
      reconnectTimer.current = undefined;
    }
    reconnectAttempt.current = 0;
  }, []);

  useEffect(() => {
    return () => {
      clearReconnect();
      if (currentBlobUrl.current) {
        URL.revokeObjectURL(currentBlobUrl.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [clearReconnect]);

  useEffect(() => {
    if (remoteEnabled) {
      containerRef.current?.focus();
    }
  }, [remoteEnabled]);

  const scheduleReconnect = useCallback(() => {
    if (isManualDisconnect.current) return;
    const delay = Math.min(1000 * Math.pow(2, reconnectAttempt.current), 30000);
    reconnectAttempt.current++;
    reconnectTimer.current = setTimeout(() => {
      if (isManualDisconnect.current) return;
      connect();
    }, delay);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const connect = useCallback(() => {
    const token = localStorage.getItem("sentineldesk_token");
    if (!token) return;

    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const wsBase = apiUrl.replace(/^http/, "ws");
    const url = `${wsBase}/ws/live/view/${deviceId}?token=${token}`;

    setStatus("connecting");
    const ws = new WebSocket(url);
    let seq = 0;

    ws.onopen = () => {
      setStatus("connected");
      reconnectAttempt.current = 0;
    };

    ws.onmessage = (event) => {
      if (event.data instanceof Blob) {
        seq++;
        const currentSeq = seq;

        frameCount.current++;

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
          setLastFrameTime(new Date().toLocaleTimeString());
          lastFpsUpdate.current = now;
        }

        if (currentSeq < lastSeq.current) return;
        lastSeq.current = currentSeq;

        if (currentBlobUrl.current) {
          URL.revokeObjectURL(currentBlobUrl.current);
        }
        const blobUrl = URL.createObjectURL(event.data);
        currentBlobUrl.current = blobUrl;
        if (imgRef.current) {
          imgRef.current.src = blobUrl;
        }
      } else if (typeof event.data === "string") {
        handleWSMessage(event.data);
      }
    };

    ws.onclose = () => {
      setStatus("disconnected");
      disableRemote();
      scheduleReconnect();
    };

    ws.onerror = () => {
      setStatus("disconnected");
    };

    wsRef.current = ws;
  }, [deviceId, handleWSMessage, disableRemote, scheduleReconnect]);

  const disconnect = useCallback(() => {
    isManualDisconnect.current = true;
    clearReconnect();
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setStatus("disconnected");
    setDisplayFps(0);
    setDisplayResolution("");
    setLastFrameTime("");
    frameTimestamps.current = [];
    frameCount.current = 0;
    lastFpsUpdate.current = 0;
    lastSeq.current = 0;
    disableRemote();
  }, [disableRemote, clearReconnect]);

  const reconnect = useCallback(() => {
    isManualDisconnect.current = false;
    clearReconnect();
    disconnect();
    setTimeout(connect, 500);
  }, [disconnect, connect, clearReconnect]);

  const handleImgLoad = () => {
    if (imgRef.current) {
      setDisplayResolution(`${imgRef.current.naturalWidth}x${imgRef.current.naturalHeight}`);
    }
  };

  const getButtonName = (button: number): string => {
    if (button === 0) return "left";
    if (button === 2) return "right";
    return "middle";
  };

  const handleMouseEvent = useCallback(
    (e: React.MouseEvent) => {
      if (!remoteEnabled) return;
      if (e.type === "mousemove") {
        sendMouseMove(e.clientX, e.clientY);
      } else if (e.type === "mousedown") {
        sendMouseDown(e.clientX, e.clientY, getButtonName(e.button));
      } else if (e.type === "mouseup") {
        sendMouseUp(e.clientX, e.clientY, getButtonName(e.button));
      }
    },
    [remoteEnabled, sendMouseMove, sendMouseDown, sendMouseUp]
  );

  const handleDoubleClickEvent = useCallback(
    (e: React.MouseEvent) => {
      if (!remoteEnabled) return;
      sendDoubleClick(e.clientX, e.clientY, getButtonName(e.button));
    },
    [remoteEnabled, sendDoubleClick]
  );

  const handleWheel = useCallback(
    (e: React.WheelEvent) => {
      if (!remoteEnabled) return;
      e.preventDefault();
      sendMouseWheel(e.clientX, e.clientY, e.deltaY);
    },
    [remoteEnabled, sendMouseWheel]
  );

  const handleKeyEvent = useCallback(
    (e: React.KeyboardEvent) => {
      if (!remoteEnabled) return;
      hookHandleKeyEvent(e);
    },
    [remoteEnabled, hookHandleKeyEvent]
  );

  const cursorStyle: React.CSSProperties | undefined = remoteEnabled
    ? {
        position: "absolute",
        left: `${(cursorX / screenInfo.width) * 100}%`,
        top: `${(cursorY / screenInfo.height) * 100}%`,
        pointerEvents: "none",
        zIndex: 50,
        transform: "translate(-50%, -50%)",
        transition: "left 0.05s linear, top 0.05s linear",
      }
    : undefined;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Monitor className="h-5 w-5" />
          Live View
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap gap-2">
          <Button
            size="sm"
            onClick={connect}
            disabled={status === "connected" || status === "connecting"}
          >
            <Play className="mr-1 h-4 w-4" />
            Start Live View
          </Button>
          <Button
            size="sm"
            variant="destructive"
            onClick={disconnect}
            disabled={status === "disconnected"}
          >
            <Square className="mr-1 h-4 w-4" />
            Stop
          </Button>
          <Button
            size="sm"
            variant="outline"
            onClick={reconnect}
            disabled={status === "connecting"}
          >
            <RefreshCw className="mr-1 h-4 w-4" />
            Reconnect
          </Button>
          {status === "connected" && (
            <Button
              size="sm"
              variant={remoteEnabled ? "default" : "secondary"}
              onClick={remoteEnabled ? disableRemote : enableRemote}
            >
              <MousePointer className="mr-1 h-4 w-4" />
              {remoteEnabled ? "Disable Remote Control" : "Enable Remote Control"}
            </Button>
          )}
        </div>

        <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
          <span>
            Status:{" "}
            {status === "connected" ? (
              <span className="text-green-500 font-medium">Connected</span>
            ) : status === "connecting" ? (
              <span className="text-yellow-500 font-medium">Connecting</span>
            ) : (
              <span className="text-red-500 font-medium">Disconnected</span>
            )}
          </span>
          <span>FPS: {displayFps}</span>
          {displayResolution && <span>Resolution: {displayResolution}</span>}
          {lastFrameTime && <span>Last Frame: {lastFrameTime}</span>}
          <span>
            Remote:{" "}
            {remoteEnabled ? (
              <Badge variant="default" className="bg-green-600">Active</Badge>
            ) : (
              <Badge variant="secondary">Off</Badge>
            )}
          </span>
          <span>Packets: {packetsSent}</span>
          {remoteEnabled && screenInfo && (
            <span>Remote Screen: {screenInfo.width}x{screenInfo.height}</span>
          )}
        </div>

        <div
          ref={containerRef}
          className="bg-muted rounded-lg overflow-hidden border relative outline-none"
          style={{ aspectRatio: "16/9" }}
          tabIndex={remoteEnabled ? 0 : -1}
          onMouseMove={handleMouseEvent}
          onMouseDown={handleMouseEvent}
          onMouseUp={handleMouseEvent}
          onDoubleClick={handleDoubleClickEvent}
          onWheel={handleWheel}
          onKeyDown={handleKeyEvent}
          onKeyUp={handleKeyEvent}
        >
          {status === "disconnected" ? (
            <div className="flex items-center justify-center h-full text-muted-foreground">
              Click &quot;Start Live View&quot; to begin
            </div>
          ) : (
            <>
              <img
                ref={imgRef}
                alt="Live Screen"
                className="w-full h-full object-contain select-none pointer-events-none"
                draggable={false}
                onLoad={handleImgLoad}
              />
              {remoteEnabled && (
                <div className="absolute inset-0 bg-transparent" />
              )}
              {remoteEnabled && (
                <div style={cursorStyle}>
                  <svg
                    width="20"
                    height="20"
                    viewBox="0 0 20 20"
                    fill="none"
                    style={{ filter: "drop-shadow(0 0 2px rgba(0,0,0,0.8))" }}
                  >
                    <path
                      d="M2 2L11 17L13 12L18 11L2 2Z"
                      fill="rgba(255,255,255,0.9)"
                      stroke="rgba(0,0,0,0.8)"
                      strokeWidth="1.5"
                    />
                  </svg>
                </div>
              )}
            </>
          )}
        </div>

        {remoteEnabled && (
          <p className="text-xs text-muted-foreground text-center">
            Remote control active. Mouse and keyboard input is being sent to the
            remote machine. Click the cursor area first to focus keyboard input.
          </p>
        )}
      </CardContent>
    </Card>
  );
}
