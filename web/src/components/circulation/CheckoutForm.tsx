import { useState } from "react";
import { checkoutBook } from "../../lib/api";

interface CheckoutFormProps {
  onSuccess?: (loan: unknown) => void;
}

export default function CheckoutForm({ onSuccess }: CheckoutFormProps) {
  const [bookId, setBookId] = useState("");
  const [patronId, setPatronId] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const loan = await checkoutBook({ book_id: bookId, patron_id: patronId });
      onSuccess?.(loan);
      setBookId("");
      setPatronId("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to checkout book");
    } finally {
      setSubmitting(false);
    }
  };

  const inputClass =
    "w-full px-3 py-2 rounded border border-book-cloth bg-paper text-ink text-body focus:outline-none focus:border-leather focus:ring-1 focus:ring-leather";

  return (
    <form onSubmit={handleSubmit} class="space-y-5 max-w-xl">
      {error && <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>}

      <div>
        <label for="bookId" class="block text-data text-reading-lamp mb-1">Book ID *</label>
        <input id="bookId" type="text" required value={bookId} onInput={(e) => setBookId((e.target as HTMLInputElement).value)} class={inputClass} placeholder="MongoDB ObjectID" />
      </div>
      <div>
        <label for="patronId" class="block text-data text-reading-lamp mb-1">Patron ID *</label>
        <input id="patronId" type="text" required value={patronId} onInput={(e) => setPatronId((e.target as HTMLInputElement).value)} class={inputClass} placeholder="MongoDB ObjectID" />
      </div>

      <button
        type="submit"
        disabled={submitting}
        class="px-6 py-2 rounded bg-leather text-paper font-body font-medium text-body hover:bg-leather/90 transition-colors disabled:opacity-50"
      >
        {submitting ? "Processing..." : "Checkout Book"}
      </button>
    </form>
  );
}
