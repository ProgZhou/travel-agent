import type { Phase } from '../types';
import './StatusBar.css';

const PHASE_LABELS: Record<Phase, string> = {
  COLLECTING: '正在收集您的出行需求...',
  RECOMMENDING: '正在为您推荐目的地...',
  PLANNING: '正在生成详细行程...',
  BOOKING: '预订引导中 - 请逐项完成预订',
  COMPLETED: '行程规划完成',
};

interface StatusBarProps {
  phase: Phase;
  isStreaming: boolean;
}

export function StatusBar({ phase, isStreaming }: StatusBarProps) {
  const isCompleted = phase === 'COMPLETED';
  return (
    <div className={`status-bar ${isCompleted ? 'status-bar--complete' : ''}`}>
      <span
        className={`status-dot ${isStreaming ? 'status-dot--active' : ''}`}
        aria-hidden="true"
      />
      <span>{PHASE_LABELS[phase]}</span>
    </div>
  );
}
