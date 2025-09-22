import { useState } from "react";
import { CATEGORY_MAP } from "../constants";

export function ScheduleItem({ event, onToggleAttend }) {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <li className={`schedule-item ${event.category}`}>
      <div className="event-detail">
        <div className="date">{event.date}</div>
        <div className="title">{event.title}</div>
        <div className="category">カテゴリ：{CATEGORY_MAP[event.category]}</div>
        
        {/* 詳細情報表示ボタン */}
        {event.details && Object.keys(event.details).length > 0 && (
          <button 
            onClick={() => setShowDetails(!showDetails)}
            style={{
              background: 'none',
              border: '1px solid #ccc',
              borderRadius: '4px',
              padding: '4px 8px',
              fontSize: '12px',
              cursor: 'pointer',
              marginTop: '8px',
              marginRight: '8px'
            }}
          >
            {showDetails ? '詳細を隠す' : '詳細を見る'}
          </button>
        )}
        
        {/* 詳細情報の表示 */}
        {showDetails && event.details && (
          <div style={{
            marginTop: '12px',
            padding: '12px',
            backgroundColor: '#f8f9fa',
            borderRadius: '4px',
            fontSize: '14px',
            border: '1px solid #e9ecef'
          }}>
            <h4 style={{ margin: '0 0 8px 0', fontSize: '14px', fontWeight: 'bold' }}>詳細情報</h4>
            {Object.entries(event.details).map(([key, value]) => (
              <div key={key} style={{ marginBottom: '4px' }}>
                <strong>{key}:</strong> {
                  // URLの場合はリンクとして表示
                  (key.includes('url') || value.startsWith('http')) ? (
                    <a 
                      href={value} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      style={{ color: '#007bff', textDecoration: 'underline' }}
                    >
                      {value}
                    </a>
                  ) : value
                }
              </div>
            ))}
          </div>
        )}
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
