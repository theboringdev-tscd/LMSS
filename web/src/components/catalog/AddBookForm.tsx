import { useState } from "react";
import { createBook, type Book } from "../../lib/api";

interface AddBookFormProps {
  onSuccess?: (book: Book) => void;
}

export default function AddBookForm({ onSuccess }: AddBookFormProps) {
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [isbn, setIsbn] = useState("");
  const [genre, setGenre] = useState("");
  const [location, setLocation] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const book = await createBook({
        title,
        author,
        isbn: isbn || undefined,
        genre: genre || undefined,
        location: location || undefined,
        status: "available",
      });
      onSuccess?.(book);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add book");
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
        <label for="title" class="block text-data text-reading-lamp mb-1">Title *</label>
        <input id="title" type="text" required value={title} onInput={(e) => setTitle((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="author" class="block text-data text-reading-lamp mb-1">Author *</label>
        <input id="author" type="text" required value={author} onInput={(e) => setAuthor((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="isbn" class="block text-data text-reading-lamp mb-1">ISBN</label>
        <input id="isbn" type="text" value={isbn} onInput={(e) => setIsbn((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="genre" class="block text-data text-reading-lamp mb-1">Genre</label>
        <input id="genre" type="text" value={genre} onInput={(e) => setGenre((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>
      <div>
        <label for="location" class="block text-data text-reading-lamp mb-1">Shelf Location</label>
        <input id="location" type="text" value={location} onInput={(e) => setLocation((e.target as HTMLInputElement).value)} class={inputClass} />
      </div>

      <button
        type="submit"
        disabled={submitting}
        class="px-6 py-2 rounded bg-leather text-paper font-body font-medium text-body hover:bg-leather/90 transition-colors disabled:opacity-50"
      >
        {submitting ? "Adding..." : "Add to Catalog"}
      </button>
    </form>
  );
}
