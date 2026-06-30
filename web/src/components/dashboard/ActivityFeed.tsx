import type { ActivityEntry } from "../../lib/api";

interface ActivityFeedProps {
  entries: ActivityEntry[];
}

export default function ActivityFeed({ entries }: ActivityFeedProps) {
  if (entries.length === 0) {
    return (
      <p className="text-reading-lamp text-body">
        No recent activity to display.
      </p>
    );
  }

  return (
    <ul className="space-y-3">
      {entries.map((entry, i) => (
        <li key={i} className="flex items-center justify-between text-body">
          <span className="text-ink">
            <em className="not-italic font-medium">{entry.book_title}</em>
            {" \u2192 "}
            <span className="text-reading-lamp">{entry.action}</span>
            {" by "}
            <span className="text-reading-lamp">{entry.patron_name}</span>
          </span>
          <span className="text-small text-reading-lamp font-mono">
            {entry.timestamp}
          </span>
        </li>
      ))}
    </ul>
  );
}
