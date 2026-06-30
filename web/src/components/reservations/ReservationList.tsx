import { useState, useEffect } from "react";
import { listReservations, deleteReservation, type Reservation } from "../../lib/api";

export interface Reservation {
  id: string;
  book_id: string;
  book_title: string;
  patron_id: string;
  patron_name: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export default function ReservationList() {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const loadReservations = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await listReservations();
      setReservations(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load reservations");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadReservations();
  }, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Cancel this reservation?")) return;
    setDeletingId(id);
    setError(null);
    try {
      await deleteReservation(id);
      setReservations((prev) => prev.filter((r) => r.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete reservation");
    } finally {
      setDeletingId(null);
    }
  };

  if (loading) {
    return (
      <div class="flex items-center justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-book-cloth border-t-leather" />
        <span class="ml-2 text-reading-lamp text-body">Loading reservations...</span>
      </div>
    );
  }
  if (error) {
    return <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>;
  }
  if (reservations.length === 0) {
    return <p class="text-reading-lamp text-body">No reservations found.</p>;
  }

  return (
    <div class="space-y-3">
      {reservations.map((res) => (
        <div key={res.id} class="flex items-center justify-between bg-book-cloth/30 rounded-card p-4 border border-book-cloth/50">
          <div>
            <h3 class="font-body font-semibold text-ink text-section">{res.book_title}</h3>
            <p class="text-reading-lamp text-data mt-1">Patron: {res.patron_name}</p>
            <p class="text-small text-reading-lamp/70 mt-1 font-mono">ID: {res.id}</p>
          </div>
          <div class="flex items-center gap-3">
            <span class={["text-small px-2 py-0.5 rounded", res.status === "pending" ? "bg-gold-leaf/20 text-binding" : "bg-green-900/20 text-green-700"].join(" ")}>
              {res.status}
            </span>
            <button
              onClick={() => handleDelete(res.id)}
              disabled={deletingId === res.id}
              class="px-3 py-1 rounded bg-red-900/20 text-red-700 text-small hover:bg-red-900/30 transition-colors disabled:opacity-50"
            >
              {deletingId === res.id ? "Canceling..." : "Cancel"}
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
