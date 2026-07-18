"use client";

import { useCallback, useRef, useState, useEffect } from "react";

interface ScreenInfo {
  width: number;
  height: number;
}

export function useRemoteControl(
  wsRef: React.MutableRefObject<WebSocket | null>,
  containerRef: React.RefObject<HTMLDivElement | null>
) {
  const [enabled, setEnabled] = useState(false);
  const [screenInfo, setScreenInfo] = useState<ScreenInfo>({ width: 1280, height: 720 });
  const [cursorX, setCursorX] = useState(0);
  const [cursorY, setCursorY] = useState(0);
  const [isConnected, setIsConnected] = useState(false);
  const [packetsSent, setPacketsSent] = useState(0);

  const packetsRef = useRef(0);
  const lastMoveRef = useRef(0);
  const moveInterval = 1000 / 60;
  const pendingMoveRef = useRef<{ x: number; y: number } | null>(null);

  const activeModifiers = useRef<Set<string>>(new Set());

  const modifierKeys = new Set([
    "ControlLeft", "ControlRight",
    "ShiftLeft", "ShiftRight",
    "AltLeft", "AltRight",
    "MetaLeft", "MetaRight",
  ]);

  const send = useCallback(
    (msg: object) => {
      if (wsRef.current?.readyState !== WebSocket.OPEN) return;
      wsRef.current.send(JSON.stringify(msg));
      packetsRef.current++;
      setPacketsSent(packetsRef.current);
    },
    [wsRef]
  );

  const scaleCoords = useCallback(
    (clientX: number, clientY: number) => {
      const rect = containerRef.current?.getBoundingClientRect();
      if (!rect) return { x: 0, y: 0 };
      const dpr = window.devicePixelRatio || 1;
      const x = Math.round(((clientX - rect.left) / rect.width) * screenInfo.width * dpr);
      const y = Math.round(((clientY - rect.top) / rect.height) * screenInfo.height * dpr);
      return { x, y };
    },
    [containerRef, screenInfo]
  );

  const sendMouseMove = useCallback(
    (clientX: number, clientY: number) => {
      const { x, y } = scaleCoords(clientX, clientY);
      const now = Date.now();
      if (now - lastMoveRef.current < moveInterval) {
        pendingMoveRef.current = { x, y };
        return;
      }
      lastMoveRef.current = now;
      send({ type: "mouse_move", x, y });
      setCursorX(x);
      setCursorY(y);
    },
    [scaleCoords, send]
  );

  const flushMove = useCallback(() => {
    if (pendingMoveRef.current) {
      const { x, y } = pendingMoveRef.current;
      pendingMoveRef.current = null;
      send({ type: "mouse_move", x, y });
      setCursorX(x);
      setCursorY(y);
    }
  }, [send]);

  const sendMouseDown = useCallback(
    (clientX: number, clientY: number, button: string) => {
      const { x, y } = scaleCoords(clientX, clientY);
      send({ type: "mouse_down", button, x, y });
    },
    [scaleCoords, send]
  );

  const sendMouseUp = useCallback(
    (clientX: number, clientY: number, button: string) => {
      const { x, y } = scaleCoords(clientX, clientY);
      send({ type: "mouse_up", button, x, y });
    },
    [scaleCoords, send]
  );

  const sendDoubleClick = useCallback(
    (clientX: number, clientY: number, button: string) => {
      const { x, y } = scaleCoords(clientX, clientY);
      send({ type: "double_click", button, x, y });
    },
    [scaleCoords, send]
  );

  const sendMouseWheel = useCallback(
    (clientX: number, clientY: number, delta: number) => {
      const { x, y } = scaleCoords(clientX, clientY);
      send({ type: "mouse_wheel", delta: Math.round(delta), x, y });
    },
    [scaleCoords, send]
  );

  const sendKeyDown = useCallback((code: string) => {
    const isModifier = modifierKeys.has(code);
    const modifiers = isModifier ? [] : Array.from(activeModifiers.current);
    send({ type: "key_down", code, modifiers });
    if (!isModifier) {
      for (const m of modifiers) {
        send({ type: "key_up", code: m, modifiers: [] });
      }
      send({ type: "key_up", code, modifiers: [] });
      for (const m of modifiers) {
        send({ type: "key_down", code: m, modifiers: [] });
      }
    }
  }, [send]);

  const sendKeyUp = useCallback((code: string) => {
    const isModifier = modifierKeys.has(code);
    const modifiers = isModifier ? [] : Array.from(activeModifiers.current);
    send({ type: "key_up", code, modifiers });
  }, [send]);

  const enable = useCallback(() => setEnabled(true), []);
  const disable = useCallback(() => {
    setEnabled(false);
    activeModifiers.current.clear();
  }, []);

  const handleKeyEvent = useCallback((e: React.KeyboardEvent) => {
    if (!enabled) return;
    e.preventDefault();
    e.stopPropagation();

    const code = e.code;
    const isModifier = modifierKeys.has(code);

    if (isModifier) {
      if (e.type === "keydown") {
        activeModifiers.current.add(code);
        send({ type: "key_down", code, modifiers: [] });
      } else {
        activeModifiers.current.delete(code);
        send({ type: "key_up", code, modifiers: [] });
      }
      return;
    }

    if (e.type === "keydown") {
      const mods = Array.from(activeModifiers.current);
      send({ type: "key_down", code, modifiers: mods });
    } else {
      const mods = Array.from(activeModifiers.current);
      send({ type: "key_up", code, modifiers: mods });
    }
  }, [enabled, send]);

  const handleWSMessage = useCallback((data: string) => {
    try {
      const msg = JSON.parse(data);
      if (msg.type === "screen_info") {
        setScreenInfo({ width: msg.width, height: msg.height });
        setIsConnected(true);
      }
    } catch {
      // ignore
    }
  }, []);

  useEffect(() => {
    const id = setInterval(() => {
      flushMove();
    }, moveInterval);
    return () => clearInterval(id);
  }, [flushMove]);

  useEffect(() => {
    const handleGlobalKeyUp = (e: KeyboardEvent) => {
      const code = e.code;
      if (modifierKeys.has(code)) {
        activeModifiers.current.delete(code);
      }
    };
    window.addEventListener("keyup", handleGlobalKeyUp);
    return () => window.removeEventListener("keyup", handleGlobalKeyUp);
  }, []);

  return {
    enabled,
    enable,
    disable,
    screenInfo,
    cursorX,
    cursorY,
    isConnected,
    packetsSent,
    sendMouseMove,
    sendMouseDown,
    sendMouseUp,
    sendDoubleClick,
    sendMouseWheel,
    sendKeyDown,
    sendKeyUp,
    handleWSMessage,
    handleKeyEvent,
  };
}
