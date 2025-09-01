import { useState } from "react";
import { eventsData } from "./demoData";

export default function App() {
  const [events, setEvents] = useState(eventsData);
  // 絞り込み状態を記憶する新しいstateを追加
  const [selectedCategory, setSelectedCategory]  = useState("all"); // useStateは第一引数（変数）と第二引数（変数を更新する関数）を返す

  const handleToggleAttend = (idToToggle) => {
    const newEvents = events.map(event => {
      if (event.id !== idToToggle) {
        return event;
      }
      return {...event, isAttending: !event.isAttending };
    });
    setEvents(newEvents);
  }

  const filteredEvents = events.filter(event => {
    if (selectedCategory === "all") {
      return true;
    }
    return event.category === selectedCategory;
  });

  return (
    <div>
      <h1>きゅるりんってしてみて スケジュール</h1>
      <NextEventDashboard events={events} />
      <div className="filter-buttons">
        <button onClick={() => setSelectedCategory("all")}>すべて</button>
        <button onClick={() => setSelectedCategory("live")}>ライブ</button>
        <button onClick={() => setSelectedCategory("media")}>メディア</button>
        <button onClick={() => setSelectedCategory("fan-meeting")}>ファンミーティング</button>
      </div>
      <ScheduleList 
        events={filteredEvents} 
        onToggleAttend={handleToggleAttend}
      />
    </div>
  )
}

{/* 最も小さい部品、一個一個のイベントデータを格納するコンポーネント */}
function ScheduleItem({ event, onToggleAttend }) {
  return (
    <li className="schedule-item">
      <div className="date">{event.date}</div>
      <div className="title">{event.title}</div>
      <div className="category">{event.category}</div>
      <button className="attend-button" onClick={() => onToggleAttend(event.id)}>
        {event.isAttending ? "❤️" : "♡"}
      </button>
    </li>
  );
}

function ScheduleList({ events, onToggleAttend }) {
  return (
    <ul className="schedule-list">
      {events.map( event => (
        <ScheduleItem
          key={event.id}
          event={event}
          onToggleAttend={onToggleAttend}
        />
      ))}
    </ul>
  );
}

function NextEventDashboard({ events }) {

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const upcomingAttendingEvents = events
  .filter(event => event.isAttending && new Date(event.date) >= today)
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