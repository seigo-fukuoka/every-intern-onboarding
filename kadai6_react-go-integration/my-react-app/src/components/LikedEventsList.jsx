export function LikedEventsList({ events }) {
  if (events.length === 0) {
    return null;
  }

  return (
    <div className="liked-events-section">
      <h3>❤️ いいね済みイベント</h3>
      <ul className="liked-list">
        {events.map((event) => (
          <li key={event.id} className="liked-item">
            <span className="liked-date">{event.date}</span>
            <span className="liked-title">{event.title}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
