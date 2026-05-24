import type { Recommendation } from '../types';
import './DestinationCard.css';

interface DestinationCardProps {
  recommendation: Recommendation;
  onSelect: (id: string) => void;
}

export function DestinationCard({ recommendation, onSelect }: DestinationCardProps) {
  const { id, destination, reason, weather, transport, hotel_range, total_estimate, sample_tickets, is_mock } = recommendation;

  return (
    <div className="dest-card" role="article" aria-label={`目的地推荐: ${destination}`}>
      <h3 className="dest-card__title">{destination}</h3>
      <p className="dest-card__reason">{reason}</p>

      <div className="dest-card__info">
        <div className="dest-card__info-item">
          <span className="label">天气</span>
          <span className="value">{weather.temp_range} {weather.condition}</span>
        </div>
        <div className="dest-card__info-item">
          <span className="label">交通</span>
          <span className="value">{transport.mode === 'flight' ? '飞机' : '高铁'} {transport.duration}</span>
        </div>
        <div className="dest-card__info-item">
          <span className="label">{sample_tickets[0]?.type === 'flight' ? '机票参考' : '高铁参考'}</span>
          <span className="value">¥{sample_tickets[0]?.price ?? transport.price}</span>
        </div>
        <div className="dest-card__info-item">
          <span className="label">酒店均价</span>
          <span className="value">¥{hotel_range.comfort}/晚</span>
        </div>
      </div>

      <div className="dest-card__total">
        预估总费用：¥{total_estimate.min.toLocaleString()} - ¥{total_estimate.max.toLocaleString()}
        {is_mock && <span className="mock-badge" title="预估数据">预估</span>}
      </div>

      <button
        className="dest-card__btn"
        onClick={() => onSelect(id)}
        aria-label={`选择 ${destination}`}
      >
        选择这个
      </button>
    </div>
  );
}
