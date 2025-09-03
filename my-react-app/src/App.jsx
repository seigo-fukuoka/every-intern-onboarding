import { useState } from "react";
import { eventsData } from "./demoData";
import './App.css'; 

const CATEGORY_MAP = {
  live: "ライブ",
  "fan-meeting": "ファンミーティング",
  media: "メディア",
  all: "すべて",
  others: "その他"
}

export default function App() {
  // useStateを使って、イベント一覧のデータを「状態」として管理する
  // events: 現在のイベント一覧データ
  // setEvents: eventsを更新するための専用関数
  const [events, setEvents] = useState(eventsData);
  // 絞り込みカテゴリを「状態」として管理する
  // selectedCategory: 現在選択されているカテゴリ名 ("all", "live" など)
  // setSelectedCategory: selectedCategoryを更新するための専用関数
  const [selectedCategory, setSelectedCategory]  = useState("all"); // useStateは第一引数（変数）と第二引数（変数を更新する関数）を返す
  const [selectedMonth, setSelectedMonth] = useState("all");

  // --- イベントハンドラ ---
  // ♡ボタンが押されたときの処理
  const handleToggleAttend = (idToToggle) => {
    // 元のevents配列を直接変更せず、新しい配列 newEvents を作成する
    
    const newEvents = events.map(event => {
      // idが合致するものがあるか照らし合わせて、なかったらそのまま返す、あればtrue or falseのisAttendingを反転させる（デフォルトはfalse）
      if (event.id !== idToToggle) {
        return event;
      }
      // 
      return {...event, isAttending: !event.isAttending };
    });
    setEvents(newEvents);
  }
  // eventsDataの「年・月」のみを取り出し、古い順にソート
  const availableMonths = [...new Set(eventsData.map(event => event.date.substring(0, 7)))].sort();  

  // 選択されたカテゴリに基づいて、表示するイベントをフィルタリングする
  const filteredEvents = events.filter(event => {
    // "すべて"が選択されている場合は、全てのイベントを対象とする
    const categoryMatch = selectedCategory ===  "all" || event.category === selectedCategory; 
    const monthMatch = selectedMonth === "all" || event.date.startsWith(selectedMonth);
    return categoryMatch && monthMatch;
    });

  const likedEvents = events.filter(event => event.isAttending);

  return (
    <div>
      <h1>きゅるりんってしてみて スケジュール</h1>
      <NextEventDashboard events={events} />
      <div className="filters">
        <select className="month-filter"
          value = {selectedMonth}
          onChange={(e) => setSelectedMonth(e.target.value)}
        >
          <option value="all">すべての月</option>
          {availableMonths.map(month => (
            <option key={month} value={month}>
              {month}
            </option>
          ))}
        </select>
        <div className="filter-buttons">
          <button onClick={() => setSelectedCategory("all")}>すべて</button>
          <button onClick={() => setSelectedCategory("live")}>ライブ</button>
          <button onClick={() => setSelectedCategory("media")}>メディア</button>
          <button onClick={() => setSelectedCategory("fan-meeting")}>ファンミーティング</button>
          <button onClick={() => setSelectedCategory("others")}>その他</button>

        </div>
      </div>
      <ScheduleList 
        events={filteredEvents} 
        onToggleAttend={handleToggleAttend}
      />
      <LikedEventsList events={likedEvents} />
    </div>
  )
}

{/* 最も小さい部品、一個一個のイベントデータを格納するコンポーネント */}
function ScheduleItem({ event, onToggleAttend }) {
  return (
    <li className="schedule-item">
      <div className="event-detail">
        <div className="date">{event.date}</div>
        <div className="title">{event.title}</div>
        <div className="category">カテゴリ：{CATEGORY_MAP[event.category]}</div>
      </div>
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

// App.jsxの末尾などに追加

function LikedEventsList({ events }) {
  // いいね済みのイベントが1つもなければ、何も表示しない
  if (events.length === 0) {
    return null;
  }

  return (
    <div className="liked-events-section">
      <h3>❤️ いいね済みイベント</h3>
      <ul className="liked-list">
        {events.map(event => (
          <li key={event.id} className="liked-item">
            {/* ↓ 日付とタイトルをspanで囲む */}
            <span className="liked-date">{event.date}</span>
            <span className="liked-title">{event.title}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}