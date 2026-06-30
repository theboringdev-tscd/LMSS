import { useState, useEffect } from "react";
import { listFines, payFine, type Fine } from "../../lib/api";
import FineList from "./FineList";

export default function FineManagement() {
  const [fines, setFines] = useState<Fine[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [payingId, setPayingId] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const loadFines = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await listFines();
      setFines(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load fines");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFines();
  }, []);

  const handlePay = async (fineId: string) => {
    if (!confirm("Mark this fine as paid?")) return;
    setPayingId(fineId);
    setError(null);
    setSuccess(null);
    try {
      await payFine(fineId);
      setSuccess("Fine marked as paid");
      setFines((prev) => prev.map((f) => (f.id === fineId ? { ...f, paid: true, status: "paid" } : f)));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to pay fine");
    } finally {
      setPayingId(null);
    }
  };

  if (loading) {
    return (
      <div class="flex items-center justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-book-cloth border-t-leather" />
        <span class="ml-2 text-reading-lamp text-body">Loading fines...</span>
      </div>
    );
  }

  return (
    <div class="space-y-4">
      {error && <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>}
      {success && <p class="text-green-700 text-body bg-green-50 rounded p-3">{success}</p>}
      <FineList fines={fines} onPay={handlePay} payingId={payingId ?? undefined} />
    </div>
  );
}
