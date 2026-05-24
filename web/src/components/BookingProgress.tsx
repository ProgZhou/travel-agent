import type { BookingItem } from '../types';
import './BookingProgress.css';

interface BookingProgressProps {
  items: BookingItem[];
  onPaid: (id: string) => void;
  onSkip: (id: string) => void;
}

export function BookingProgress({ items, onPaid, onSkip }: BookingProgressProps) {
  const sortedItems = [...items].sort((a, b) => a.order - b.order);
  const doneCount = sortedItems.filter((i) => i.status !== 'PENDING').length;
  const progress = (doneCount / sortedItems.length) * 100;
  const currentItem = sortedItems.find((i) => i.status === 'PENDING');

  return (
    <div role="region" aria-label="预订进度">
      <div className="booking-progress">
        <div className="progress-bar" role="progressbar" aria-valuenow={progress} aria-valuemin={0} aria-valuemax={100}>
          <div className="progress-fill" style={{ width: `${progress}%` }} />
        </div>
        <p className="progress-text">{doneCount}/{sortedItems.length} 已完成</p>
      </div>

      {sortedItems.map((item) => {
        const isCurrent = item.id === currentItem?.id;
        const isDone = item.status === 'BOOKED';
        const isSkipped = item.status === 'SKIPPED';

        return (
          <div
            key={item.id}
            className={`booking-item ${isCurrent ? 'booking-item--current' : ''} ${isDone ? 'booking-item--done' : ''} ${isSkipped ? 'booking-item--skipped' : ''}`}
          >
            <h4>
              {isDone && <span aria-hidden="true">✓</span>}
              {item.item_name}
            </h4>
            <p className="detail">{item.detail}</p>
            <span
              className={`status-badge ${item.status === 'PENDING' ? 'status-badge--pending' : item.status === 'BOOKED' ? 'status-badge--done' : 'status-badge--skipped'}`}
            >
              {item.status === 'PENDING' ? '待预订' : item.status === 'BOOKED' ? '已预订' : '已跳过'}
            </span>

            {isCurrent && (
              <div className="booking-actions">
                <button
                  className="btn-book"
                  onClick={() => window.open(item.link, '_blank')}
                  aria-label={`去${item.platform}预订`}
                >
                  去{item.platform}预订
                </button>
                <button
                  className="btn-paid"
                  onClick={() => onPaid(item.id)}
                  aria-label="标记为已支付"
                >
                  已支付
                </button>
                <button
                  className="btn-skip"
                  onClick={() => onSkip(item.id)}
                  aria-label="跳过此项"
                >
                  跳过
                </button>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
