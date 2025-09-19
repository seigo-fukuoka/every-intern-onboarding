export function NextEventDashboard({ events }) {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const upcomingAttendingEvents = events
    .filter((event) => event.isAttending && new Date(event.date) >= today)
    .sort((a, b) => new Date(a.date) - new Date(b.date));

  const nextEvent = upcomingAttendingEvents[0];

  let message = "参加予定のイベントはありません";

  if (nextEvent) {
    const nextEventDate = new Date(nextEvent.date);
    const diffTime = nextEventDate - today;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24) - 1);

    message = `次の予定「${nextEvent.title}」まであと${diffDays}日！`;
  }
  return (
    <div className="dashboard">
      <h2>{message}</h2>
    </div>
  );
}
