import { useState, useEffect } from "react";
import { getPatron, deletePatron, type Patron } from "../../lib/api";

interface PatronDetailProps {
  bookId?: string;
}

export default function PatronDetail({ bookId }: PatronDetailProps) {
  const [patron, setPatron] = useState<Patron | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const id = bookId || window.location.pathname.split("/").pop();

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    setError(null);
    getPatron(id)
      .then(setPatron)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load patron"))
      .finally(() => setLoading(false));
  }, [id]);

  const handleDelete = async () => {
    if (!id || !confirm("Are you sure you want to delete this patron?")) return;
    try {
      await deletePatron(id);
      window.location.href = "/patrons";
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete patron");
    }
  };

  if (loading) {
    return <p class="text-reading-lamp text-body">Loading patron...</p>;
  }
  if (error) {
    return <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>;
  }
  if (!patron) {
    return <p class="text-reading-lamp text-body">Patron not found.</p>;
  }

  return (
    <div class="space-y-6">
      <div>
        <h2 class="font-display text-page-title text-ink">{patron.name}</h2>
        <p class="text-reading-lamp text-data mt-1">{patron.email}</p>
      </div>
      <hr class="gold-rule" aria-hidden="true" />
      <dl class="grid grid-cols-1 gap-4 max-w-xl">
        <div>
          <dt class="text-data text-reading-lamp">Email</dt>
          <dd class="text-ink text-body mt-1">{patron.email}</dd>
        </div>
        {patron.phone && (
          <div>
            <dt class="text-data text-reading-lamp">Phone</dt>
            <dd class="text-ink text-body mt-1 font-mono">{patron.phone}</dd>
          </div>
        )}
        <div>
          <dt class="text-data text-reading-lamp">Status</dt>
          <dd class="mt-1">
            <span
              class={[
                "text-small px-2 py-0.5 rounded",
                patron.active ? "bg-green-900/20 text-green-700" : "bg-red-900/20 text-red-700",
              ].join(" ")}
            >
              {patron.active ? "Active" : "Inactive"}
            </span>
          </dd>
        </div>
        <div>
          <dt class="text-data text-reading-lamp">Added</dt>
          <dd class="text-ink text-body mt-1">{new Date(patron.created_at).toLocaleDateString()}</dd>
        </div>
      </dl>
      <button
        onClick={handleDelete}
        class="px-4 py-2 rounded bg-red-900/20 text-red-700 text-body hover:bg-red-900/30 transition-colors"
      >
        Delete Patron
      </button>
    </div>
  );
}
