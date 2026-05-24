import { useCallback } from 'react';
import { useChatStore } from './stores/chatStore';
import { useSession } from './hooks/useSession';
import { useSSE } from './hooks/useSSE';
import { StatusBar } from './components/StatusBar';
import { ChatView } from './components/ChatView';
import { InputBar } from './components/InputBar';
import './App.css';

export function App() {
  const { sessionId } = useSession();
  const {
    phase,
    messages,
    isStreaming,
    setPhase,
    setStreaming,
    appendUserMessage,
    startAgentMessage,
    appendTextDelta,
    appendBlock,
    updateToolCallStatus,
  } = useChatStore();

  const { connect } = useSSE();

  const sendMessage = useCallback(
    async (text: string) => {
      if (!sessionId || isStreaming) return;

      appendUserMessage(text);
      setStreaming(true);

      const agentMsgId = startAgentMessage();

      await connect(
        '/api/chat',
        { session_id: sessionId, message: text },
        {
          onThinking: (data) => {
            appendBlock(agentMsgId, { type: 'thinking', content: data.content });
          },
          onText: (data) => {
            if (data.delta) {
              appendTextDelta(agentMsgId, data.content);
            } else {
              appendBlock(agentMsgId, { type: 'text', content: data.content });
            }
          },
          onToolCall: (data) => {
            if (data.status === 'running') {
              appendBlock(agentMsgId, {
                type: 'tool_call',
                name: data.name,
                status: 'running',
              });
            } else {
              updateToolCallStatus(agentMsgId, data.name, data.status);
            }
          },
          onCard: (data) => {
            appendBlock(agentMsgId, {
              type: 'card',
              card_type: data.card_type,
              data: data.data,
            });
          },
          onPhaseChange: (data) => {
            setPhase(data.to);
          },
          onDone: () => {
            setStreaming(false);
          },
          onError: (data) => {
            appendBlock(agentMsgId, {
              type: 'text',
              content: `❌ ${data.message}`,
            });
            setStreaming(false);
          },
        },
      );
    },
    [
      sessionId,
      isStreaming,
      appendUserMessage,
      setStreaming,
      startAgentMessage,
      connect,
      appendBlock,
      appendTextDelta,
      updateToolCallStatus,
      setPhase,
    ],
  );

  const handleSelectDestination = useCallback(
    (id: string) => {
      sendMessage(`选择 ${id}`);
    },
    [sendMessage],
  );

  const handleConfirmItinerary = useCallback(() => {
    sendMessage('确认行程');
  }, [sendMessage]);

  const handleModifyItinerary = useCallback(() => {
    sendMessage('我要修改行程');
  }, [sendMessage]);

  const handleBookingPaid = useCallback(
    (id: string) => {
      sendMessage(`已支付 ${id}`);
    },
    [sendMessage],
  );

  const handleBookingSkip = useCallback(
    (id: string) => {
      sendMessage(`跳过 ${id}`);
    },
    [sendMessage],
  );

  const handleComplete = useCallback(() => {
    // User can start a new session or close the app
    sendMessage('谢谢');
  }, [sendMessage]);

  if (!sessionId) {
    return (
      <div className="app-loading">
        <div className="spinner" />
        <p>正在初始化...</p>
      </div>
    );
  }

  return (
    <div className="app">
      <StatusBar phase={phase} isStreaming={isStreaming} />
      <ChatView
        messages={messages}
        isStreaming={isStreaming}
        onSelectDestination={handleSelectDestination}
        onConfirmItinerary={handleConfirmItinerary}
        onModifyItinerary={handleModifyItinerary}
        onBookingPaid={handleBookingPaid}
        onBookingSkip={handleBookingSkip}
        onComplete={handleComplete}
      />
      <InputBar onSend={sendMessage} disabled={isStreaming} />
    </div>
  );
}
