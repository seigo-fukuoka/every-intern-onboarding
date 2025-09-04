import { useState } from "react";
import { eventsData } from "./demoData";
import { LikedEventsList } from "./components/LikedEventsList";
import { NextEventDashboard } from "./components/NextEventDashboard";
import { ScheduleList } from "./components/ScheduleList";
import "./App.css";

export default function App() {
  // useStateを使って、イベント一覧のデータを「状態」として管理する
  // events: 現在のイベント一覧データ
  // setEvents: eventsを更新するための専用関数
  const [events, setEvents] = useState(eventsData);
  // 絞り込みカテゴリを「状態」として管理する
  // selectedCategory: 現在選択されているカテゴリ名 ("all", "live" など)
  // setSelectedCategory: selectedCategoryを更新するための専用関数
  const [selectedCategory, setSelectedCategory] = useState("all"); // useStateは第一引数（変数）と第二引数（変数を更新する関数）を返す
  const [selectedMonth, setSelectedMonth] = useState("all");

  // --- イベントハンドラ ---
  // ♡ボタンが押されたときの処理
  const handleToggleAttend = (idToToggle) => {
    // 元のevents配列を直接変更せず、新しい配列 newEvents を作成する

    const newEvents = events.map((event) => {
      // idが合致するものがあるか照らし合わせて、なかったらそのまま返す、あればtrue or falseのisAttendingを反転させる（デフォルトはfalse）
      if (event.id !== idToToggle) {
        return event;
      }
      //
      return { ...event, isAttending: !event.isAttending };
    });
    setEvents(newEvents);
  };
  // eventsDataの「年・月」のみを取り出し、古い順にソート
  const availableMonths = [
    ...new Set(eventsData.map((event) => event.date.substring(0, 7))),
  ].sort();

  // 選択されたカテゴリに基づいて、表示するイベントをフィルタリングする
  const filteredEvents = events.filter((event) => {
    const categoryMatch =
      selectedCategory === "all" || event.category === selectedCategory;
    const monthMatch =
      selectedMonth === "all" || event.date.startsWith(selectedMonth);
    return categoryMatch && monthMatch;
  });

  const likedEvents = events.filter((event) => event.isAttending);

  return (
    <div>
      <h1>きゅるりんってしてみて スケジュール</h1>
      <NextEventDashboard events={events} />
      <div className="filters">
        <select
          className="month-filter"
          value={selectedMonth}
          onChange={(e) => setSelectedMonth(e.target.value)}
        >
          <option value="all">すべての月</option>
          {availableMonths.map((month) => (
            <option key={month} value={month}>
              {month}
            </option>
          ))}
        </select>
        <div className="filter-buttons">
          <button onClick={() => setSelectedCategory("all")}>すべて</button>
          <button onClick={() => setSelectedCategory("live")}>ライブ</button>
          <button onClick={() => setSelectedCategory("media")}>メディア</button>
          <button onClick={() => setSelectedCategory("fanMeeting")}>
            ファンミーティング
          </button>
          <button onClick={() => setSelectedCategory("others")}>その他</button>
        </div>
      </div>
      <LikedEventsList events={likedEvents} />
      <ScheduleList
        events={filteredEvents}
        onToggleAttend={handleToggleAttend}
      />
    </div>
  );
}
