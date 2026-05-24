import type { Itinerary } from '../types';
import './ItineraryTimeline.css';

interface ItineraryTimelineProps {
  itinerary: Itinerary;
  onConfirm: () => void;
  onModify: () => void;
}

export function ItineraryTimeline({ itinerary, onConfirm, onModify }: ItineraryTimelineProps) {
  const { destination, outbound, hotel, daily_plans, return: returnFlight, cost_summary, total_cost } = itinerary;

  return (
    <div role="article" aria-label={`${destination} 行程表`}>
      <div className="itinerary-timeline">
        {/* Outbound */}
        <div className="timeline-item timeline-item--transport">
          <div className="timeline-card">
            <h4>去程交通</h4>
            <p>{outbound.number} {outbound.depart_city} → {outbound.arrive_city}</p>
            <p>{outbound.depart_time.split(' ')[0]} {outbound.depart_time.split(' ')[1]} - {outbound.arrive_time.split(' ')[1]}（{outbound.duration}）</p>
            <p className="price">¥{outbound.price} {outbound.seat}</p>
          </div>
        </div>

        {/* Hotel */}
        <div className="timeline-item timeline-item--hotel">
          <div className="timeline-card">
            <h4>住宿</h4>
            <p>{hotel.name}</p>
            <p>{hotel.room_type} · {hotel.nights}晚（{hotel.check_in} - {hotel.check_out}）</p>
            <p>{hotel.location}</p>
            <p className="price">¥{hotel.price_per_night}/晚 × {hotel.nights} = ¥{hotel.total_price}</p>
          </div>
        </div>

        {/* Daily Plans */}
        {daily_plans.map((day) => (
          <div key={day.date} className="timeline-item">
            <div className="timeline-card">
              <div className="day-title">Day {day.day_number} · {day.date} - {day.title}</div>
              <ul className="activity-list">
                {day.activities.map((act, i) => (
                  <li key={i}>
                    <span>{act.time} {act.name}</span>
                    {act.cost > 0 && <span>¥{act.cost}</span>}
                    {act.cost === 0 && <span>免费</span>}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        ))}

        {/* Return */}
        <div className="timeline-item timeline-item--transport">
          <div className="timeline-card">
            <h4>返程交通</h4>
            <p>{returnFlight.number} {returnFlight.depart_city} → {returnFlight.arrive_city}</p>
            <p>{returnFlight.depart_time.split(' ')[0]} {returnFlight.depart_time.split(' ')[1]} - {returnFlight.arrive_time.split(' ')[1]}（{returnFlight.duration}）</p>
            <p className="price">¥{returnFlight.price} {returnFlight.seat}</p>
          </div>
        </div>
      </div>

      {/* Cost Summary */}
      <div className="cost-summary">
        <h4>费用汇总</h4>
        <div className="cost-row"><span>去程交通</span><span>¥{cost_summary.outbound_cost}</span></div>
        <div className="cost-row"><span>返程交通</span><span>¥{cost_summary.return_cost}</span></div>
        <div className="cost-row"><span>住宿（{hotel.nights}晚）</span><span>¥{cost_summary.hotel_cost}</span></div>
        <div className="cost-row"><span>景点门票</span><span>¥{cost_summary.activity_cost}</span></div>
        <div className="cost-row"><span>餐饮交通</span><span>¥{cost_summary.meal_cost + cost_summary.transport_cost}</span></div>
        <div className="cost-row cost-row--total"><span>合计</span><span>¥{total_cost}</span></div>
      </div>

      {/* Action Buttons */}
      <div className="action-buttons">
        <button className="btn-secondary" onClick={onModify}>我要修改</button>
        <button className="btn-primary" onClick={onConfirm}>确认行程</button>
      </div>
    </div>
  );
}
