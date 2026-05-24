import { useState, type FormEvent } from 'react';
import './InputBar.css';

interface InputBarProps {
  onSend: (text: string) => void;
  disabled: boolean;
}

export function InputBar({ onSend, disabled }: InputBarProps) {
  const [text, setText] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = text.trim();
    if (!trimmed || disabled) return;
    onSend(trimmed);
    setText('');
  };

  return (
    <form className="input-bar" onSubmit={handleSubmit}>
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="输入你的需求..."
        disabled={disabled}
        aria-label="消息输入框"
      />
      <button type="submit" disabled={disabled || !text.trim()}>
        发送
      </button>
    </form>
  );
}
