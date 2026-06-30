import { useState, useEffect } from "react";
import StatBlock from "./StatBlock";
import ActivityFeed from "./ActivityFeed";
import { getDashboardOverview, logout } from "../../lib/api";
import type { OverviewStats } from "../../lib/api";

export default function Dashboard() {
  const [stats, setStats] = useState<OverviewStats | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getDashboardOverview()
      .then(setStats)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load dashboard"));
  }, []);

  if (error) {
    return (
      <section className="space-y-8">
        <header>
          <h1 className="font-display text-page-title text-ink">Dashboard</h1>
          <p className="text-reading-lamp text-data mt-1">Library overview for today</p>
        </header>
        <div className="rounded-card border border-book-cloth/50 bg-book-cloth/20 p-8 text-center">
          <p className="text-reading-lamp text-body">{error}</p>
        </div>
      </section>
    );
  }

  return (
    <section className="space-y-8">
      <header className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-page-title text-ink">Dashboard</h1>
          <p className="text-reading-lamp text-data mt-1">Library overview for today</p>
        </div>
        <button
          onClick={() => logout().then(() => { window.location.href = "/login"; })}
          className="text-data text-reading-lamp hover:text-leather transition-colors"
        >
          Sign out
        </button>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <StatBlock
          value={stats?.books_checked_out_today ?? "\u2014"}
          label="books checked out today"
        />
        <StatBlock
          value={stats?.new_patrons_this_week ?? "\u2014"}
          label="new patrons this week"
        />
        <StatBlock
          value={stats?.overdue_returns ?? "\u2014"}
          label="overdue returns"
        />
      </div>

      <hr className="gold-rule" aria-hidden="true" />

      <section>
        <h2 className="font-body text-section text-ink mb-4">Recent Activity</h2>
        <ActivityFeed entries={stats?.recent_activity ?? []} />
      </section>
    </section>
  );
}
