import type { BookingSummary } from '../types';
import './CompleteSummary.css';

interface CompleteSummaryProps {
  summary: BookingSummary;
  onComplete: () => void;
}

export function CompleteSummary({ summary, onComplete }: CompleteSummaryProps) {
  return (
    <div className="complete-container" role="region" aria-label="预订完成汇总">
      <div className="complete-icon" aria-hidden="true">✈️</div>
      <h2 className="complete-title">{summary.destination}之旅 预订完成!</h2>
      <p className="complete-subtitle">以下是你的预订汇总，祝旅途愉快</p>

      <div className="summary-list">
        {summary.items.map((item) => (
          <div key={item.id} className="summary-item">
            <span className="name">
              {item.status === 'BOOKED' && <span className="check-icon" aria-hidden="true">✓</span>}
              {item.status === 'SKIPPED' && <span className="skip-icon" aria-hidden="true">-</span>}
              {item.item_name}
            </span>
            <span className={`cost ${item.status === 'SKIPPED' ? 'cost--skipped' : ''}`}>
              {item.status === 'BOOKED' ? `¥${item.actual_cost || item.price}` : '已跳过'}
            </span>
          </div>
        ))}
      </div>

      <div className="total-box">
        <div className="label">实际已预订总花费</div>
        <div className="amount">¥{summary.total_booked.toLocaleString()}</div>
      </div>

      <button className="btn-complete" onClick={onComplete}>
        完成
      </button>
    </div>
  );
}
