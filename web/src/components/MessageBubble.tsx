import type { ContentBlock } from '../types';
import './MessageBubble.css';

const TOOL_NAMES: Record<string, string> = {
  flight_search: '搜索航班/火车',
  hotel_search: '搜索酒店',
  weather_query: '查询天气',
};

interface MessageBubbleProps {
  role: 'user' | 'agent';
  blocks: ContentBlock[];
  isStreaming?: boolean;
}

export function MessageBubble({ role, blocks, isStreaming }: MessageBubbleProps) {
  if (role === 'user') {
    const text = blocks.find((b) => b.type === 'text');
    return (
      <div className="message-bubble message-bubble--user" role="article" aria-label="用户消息">
        {text?.type === 'text' ? text.content : ''}
      </div>
    );
  }

  // Agent message: render each block
  return (
    <>
      {blocks.map((block, i) => {
        if (block.type === 'thinking') {
          return (
            <div key={i} className="thinking-block" role="note" aria-label="Agent 思考过程">
              <span className="thinking-icon" aria-hidden="true">💭</span>
              <span>{block.content}</span>
            </div>
          );
        }

        if (block.type === 'tool_call') {
          const label = TOOL_NAMES[block.name] ?? block.name;
          const isDone = block.status === 'done';
          const isError = block.status === 'error';
          return (
            <div
              key={i}
              className={`tool-call-block ${isDone ? 'tool-call-block--done' : ''} ${isError ? 'tool-call-block--error' : ''}`}
              role="status"
              aria-label={`工具调用: ${label}`}
            >
              {block.status === 'running' && (
                <span className="tool-spinner" aria-hidden="true" />
              )}
              {isDone && <span aria-hidden="true">✓</span>}
              {isError && <span aria-hidden="true">✗</span>}
              <span>{label}{isDone ? ' 完成' : isError ? ' 失败' : '...'}</span>
            </div>
          );
        }

        if (block.type === 'text') {
          const isLast = i === blocks.length - 1;
          return (
            <div
              key={i}
              className="message-bubble message-bubble--agent"
              role="article"
              aria-label="Agent 回复"
            >
              {block.content}
              {isStreaming && isLast && (
                <span className="cursor-blink" aria-hidden="true" />
              )}
            </div>
          );
        }

        // card blocks are rendered by the parent (ChatView)
        return null;
      })}
    </>
  );
}
