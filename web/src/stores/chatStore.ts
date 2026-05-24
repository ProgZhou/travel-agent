import { create } from 'zustand';
import type { ChatMessage, ContentBlock, Phase } from '../types';

let msgIdCounter = 0;
function nextId() {
  return `msg_${++msgIdCounter}_${Date.now()}`;
}

interface ChatStore {
  // Session state
  sessionId: string | null;
  phase: Phase;

  // Message list
  messages: ChatMessage[];

  // SSE connection state
  isConnecting: boolean;
  isStreaming: boolean;

  // Actions
  setSessionId: (id: string) => void;
  setPhase: (phase: Phase) => void;
  setConnecting: (v: boolean) => void;
  setStreaming: (v: boolean) => void;

  appendUserMessage: (content: string) => void;
  startAgentMessage: () => string; // returns message id
  appendTextDelta: (msgId: string, delta: string) => void;
  appendBlock: (msgId: string, block: ContentBlock) => void;
  updateToolCallStatus: (
    msgId: string,
    name: string,
    status: 'running' | 'done' | 'error',
  ) => void;
  reset: () => void;
}

const initialState = {
  sessionId: null as string | null,
  phase: 'COLLECTING' as Phase,
  messages: [] as ChatMessage[],
  isConnecting: false,
  isStreaming: false,
};

export const useChatStore = create<ChatStore>((set) => ({
  ...initialState,

  setSessionId: (id) => set({ sessionId: id }),
  setPhase: (phase) => set({ phase }),
  setConnecting: (v) => set({ isConnecting: v }),
  setStreaming: (v) => set({ isStreaming: v }),

  appendUserMessage: (content) =>
    set((state) => ({
      messages: [
        ...state.messages,
        {
          id: nextId(),
          role: 'user',
          blocks: [{ type: 'text', content }],
          timestamp: Date.now(),
        },
      ],
    })),

  startAgentMessage: () => {
    const id = nextId();
    set((state) => ({
      messages: [
        ...state.messages,
        {
          id,
          role: 'agent',
          blocks: [],
          timestamp: Date.now(),
        },
      ],
    }));
    return id;
  },

  appendTextDelta: (msgId, delta) =>
    set((state) => ({
      messages: state.messages.map((msg) => {
        if (msg.id !== msgId) return msg;
        const blocks = [...msg.blocks];
        const last = blocks[blocks.length - 1];
        if (last && last.type === 'text') {
          // Append to existing text block
          return {
            ...msg,
            blocks: [
              ...blocks.slice(0, -1),
              { type: 'text' as const, content: last.content + delta },
            ],
          };
        }
        // Start a new text block
        return {
          ...msg,
          blocks: [...blocks, { type: 'text' as const, content: delta }],
        };
      }),
    })),

  appendBlock: (msgId, block) =>
    set((state) => ({
      messages: state.messages.map((msg) =>
        msg.id === msgId ? { ...msg, blocks: [...msg.blocks, block] } : msg,
      ),
    })),

  updateToolCallStatus: (msgId, name, status) =>
    set((state) => ({
      messages: state.messages.map((msg) => {
        if (msg.id !== msgId) return msg;
        const blocks = msg.blocks.map((b) =>
          b.type === 'tool_call' && b.name === name && b.status === 'running'
            ? { ...b, status }
            : b,
        );
        return { ...msg, blocks };
      }),
    })),

  reset: () => set(initialState),
}));
