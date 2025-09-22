import { CATEGORY_MAP } from "../constants";

export function ScheduleItem({ event, onToggleAttend }) {
  return (
    <li className={`schedule-item ${event.category}`}>
      <div className="event-detail">
        <div className="date">{event.date}</div>
        <div className="title">{event.title}</div>
        <div className="category">カテゴリ：{CATEGORY_MAP[event.category]}</div>
      </div>
      <button
        className="attend-button"
        onClick={() => onToggleAttend(event.id)}
      >
        {event.isAttending ? "❤️" : "♡"}
      </button>
    </li>
  );
}
