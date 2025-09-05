import { ScheduleItem } from "./ScheduleItem";

export function ScheduleList({ events, onToggleAttend }) {
  return (
    <ul className="schedule-list">
      {events.map((event) => (
        <ScheduleItem
          key={event.id}
          event={event}
          onToggleAttend={onToggleAttend}
        />
      ))}
    </ul>
  );
}
