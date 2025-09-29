import { useState, useEffect } from "react";
import { LikedEventsList } from "./components/LikedEventsList.jsx";
import { NextEventDashboard } from "./components/NextEventDashboard.jsx";
import { ScheduleList } from "./components/ScheduleList.jsx";
import { FILTER_TYPE_ALL } from "./constants.js";
import "./App.css";

export default function App() {
  const [events, setEvents] = useState([]);
  const [selectedMonth, setSelectedMonth] = useState(FILTER_TYPE_ALL);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const response = await fetch("http://localhost:1323/events");
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const eventsData = await response.json();
        setEvents(eventsData);
      } catch (error) {
        setError(error.message); // エラーメッセージをセット
        console.error("Failed to fetch events:", error);
      } 
    };

    fetchEvents();
  }, []);

  // --- イベントハンドラ ---
  const handleToggleAttend = (idToToggle) => {
    const newEvents = events.map((event) => {
      if (event.id !== idToToggle) {
        return event;
      }
      const newEvent = JSON.parse(JSON.stringify(event));
      newEvent.isAttending = !newEvent.isAttending;
      return newEvent;
    });
    setEvents(newEvents);
  };


  const availableMonths = [
    ...new Set(events.map((event) => event.date.substring(0, 7))),
  ].sort();

  const filteredEvents = events.filter((event) => {
    if (
      selectedMonth !== FILTER_TYPE_ALL &&
      !event.date.startsWith(selectedMonth)
    ) {
      return false;
    }
    return true;
  });

  const likedEvents = events.filter((event) => event.isAttending);

  return (
    <div>
      <h1>きゅるりんってしてみて スケジュール</h1>
      {error && <p style={{ color: "red" }}>Error: {error}</p>} {/* エラーメッセージを表示 */}
      <NextEventDashboard events={events} />
      <div className="filters">
        <select
          className="month-filter"
          value={selectedMonth}
          onChange={(e) => setSelectedMonth(e.target.value)}
        >
          <option value={FILTER_TYPE_ALL}>すべての月</option>
          {availableMonths.map((month) => (
            <option key={month} value={month}>
              {month}
            </option>
          ))}
        </select>
      </div>
      <LikedEventsList events={likedEvents} />
      <ScheduleList
        events={filteredEvents}
        onToggleAttend={handleToggleAttend}
      />
    </div>
  );
}
