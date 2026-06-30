import { useState } from "react";
import { createPatron, type Patron } from "../../lib/api";

interface PatronFormProps {
  onSuccess?: (patron: Patron) => void;
}

export default function PatronForm({ onSuccess }: PatronFormProps) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const patron = await createPatron({
        name,
        email,
        phone: phone || undefined,
        active: true,
      });
      onSuccess?.(patron);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add patron");
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
        <label for="name" class="block text-data text-reading-lamp mb-1">Name *</label>
        <input id="name" type="text" required value={name} onInput={(e) => setName((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="email" class="block text-data text-reading-lamp mb-1">Email *</label>
        <input id="email" type="email" required value={email} onInput={(e) => setEmail((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="phone" class="block text-data text-reading-lamp mb-1">Phone</label>
        <input id="phone" type="tel" value={phone} onInput={(e) => setPhone((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>

      <button
        type="submit"
        disabled={submitting}
        class="px-6 py-2 rounded bg-leather text-paper font-body font-medium text-body hover:bg-leather/90 transition-colors disabled:opacity-50"
      >
        {submitting ? "Adding..." : "Add Patron"}
      </button>
    </form>
  );
}
