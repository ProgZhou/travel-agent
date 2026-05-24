import { useCallback, useRef } from 'react';
import type {
  CardEventData,
  ErrorEventData,
  PhaseChangeEventData,
  TextEventData,
  ThinkingEventData,
  ToolCallEventData,
} from '../types';

export interface SSEHandlers {
  onThinking?: (data: ThinkingEventData) => void;
  onText?: (data: TextEventData) => void;
  onToolCall?: (data: ToolCallEventData) => void;
  onCard?: (data: CardEventData) => void;
  onPhaseChange?: (data: PhaseChangeEventData) => void;
  onDone?: () => void;
  onError?: (data: ErrorEventData) => void;
}

export function useSSE() {
  const abortRef = useRef<AbortController | null>(null);

  const connect = useCallback(
    async (
      url: string,
      body: Record<string, unknown>,
      handlers: SSEHandlers,
    ) => {
      // Cancel any existing connection
      abortRef.current?.abort();
      const controller = new AbortController();
      abortRef.current = controller;

      try {
        const res = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'text/event-stream',
          },
          body: JSON.stringify(body),
          signal: controller.signal,
        });

        if (!res.ok) {
          // Non-SSE error response
          const json = await res.json().catch(() => ({ message: 'Request failed' }));
          handlers.onError?.({
            code: String(res.status),
            message: (json as { message?: string }).message ?? 'Request failed',
          });
          return;
        }

        const reader = res.body?.getReader();
        if (!reader) {
          handlers.onError?.({ code: '0', message: 'No response body' });
          return;
        }

        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });

          // SSE messages are separated by double newlines
          const parts = buffer.split('\n\n');
          buffer = parts.pop() ?? '';

          for (const part of parts) {
            if (!part.trim()) continue;

            let eventType = 'message';
            let dataLine = '';

            for (const line of part.split('\n')) {
              if (line.startsWith('event: ')) {
                eventType = line.slice(7).trim();
              } else if (line.startsWith('data: ')) {
                dataLine = line.slice(6).trim();
              }
            }

            if (!dataLine) continue;

            let parsed: unknown;
            try {
              parsed = JSON.parse(dataLine);
            } catch {
              continue;
            }

            switch (eventType) {
              case 'thinking':
                handlers.onThinking?.(parsed as ThinkingEventData);
                break;
              case 'text':
                handlers.onText?.(parsed as TextEventData);
                break;
              case 'tool_call':
                handlers.onToolCall?.(parsed as ToolCallEventData);
                break;
              case 'card':
                handlers.onCard?.(parsed as CardEventData);
                break;
              case 'phase_change':
                handlers.onPhaseChange?.(parsed as PhaseChangeEventData);
                break;
              case 'done':
                handlers.onDone?.();
                break;
              case 'error':
                handlers.onError?.(parsed as ErrorEventData);
                break;
            }
          }
        }
      } catch (err) {
        if ((err as Error).name === 'AbortError') return;
        handlers.onError?.({
          code: '0',
          message: '连接失败，请检查网络后重试',
        });
      }
    },
    [],
  );

  const disconnect = useCallback(() => {
    abortRef.current?.abort();
    abortRef.current = null;
  }, []);

  return { connect, disconnect };
}
