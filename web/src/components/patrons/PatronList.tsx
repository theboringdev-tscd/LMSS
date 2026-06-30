import { useState, useEffect } from "react";
import { listPatrons, searchPatrons, type Patron } from "../../lib/api";
import PatronCard from "./PatronCard";

export default function PatronList() {
  const [patrons, setPatrons] = useState<Patron[]>([]);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadPatrons = async (q: string) => {
    setLoading(true);
    setError(null);
    try {
      const data = q ? await searchPatrons(q) : await listPatrons();
      setPatrons(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load patrons");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPatrons("");
  }, []);

  const handleSearch = (q: string) => {
    setQuery(q);
    loadPatrons(q);
  };

  return (
    <div class="space-y-4">
      <div class="flex items-center gap-3">
        <input
          type="text"
          value={query}
          onInput={(e) => handleSearch((e.target as HTMLInputElement).value)}
          placeholder="Search by name or email..."
          class="flex-1 px-3 py-2 rounded border border-book-cloth bg-paper text-ink text-body focus:outline-none focus:border-leather focus:ring-1 focus:ring-leather"
        />
        {query && (
          <button
            onClick={() => handleSearch("")}
            class="px-3 py-2 rounded bg-book-cloth text-ink text-body hover:bg-book-cloth/80 transition-colors"
          >
            Clear
          </button>
        )}
      </div>
      {error && <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>}
      {loading ? (
        <div class="flex items-center justify-center py-8">
          <div class="h-6 w-6 animate-spin rounded-full border-2 border-book-cloth border-t-leather" />
          <span class="ml-2 text-reading-lamp text-body">Loading patrons...</span>
        </div>
      ) : patrons.length === 0 ? (
        <p class="text-reading-lamp text-body">No patrons found.</p>
      ) : (
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {patrons.map((patron) => (
            <PatronCard key={patron.id} patron={patron} />
          ))}
        </div>
      )}
    </div>
  );
}
