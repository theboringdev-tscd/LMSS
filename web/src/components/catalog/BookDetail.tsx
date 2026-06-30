import { useState, useEffect } from "react";
import { getBook, deleteBook } from "../../lib/api";
import type { Book } from "./BookCard";

interface BookDetailProps {
  bookId: string;
}

export default function BookDetail({ bookId }: BookDetailProps) {
  const [book, setBook] = useState<Book | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    getBook(bookId)
      .then(setBook)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load book"));
  }, [bookId]);

  const handleDelete = async () => {
    if (!confirm(`Remove "${book?.title}" from the catalog?`)) return;
    setDeleting(true);
    try {
      await deleteBook(bookId);
      window.location.href = "/catalog";
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete book");
      setDeleting(false);
    }
  };

  if (error) {
    return <p class="text-red-700 text-body bg-red-50 rounded p-4">{error}</p>;
  }

  if (!book) {
    return <p class="text-reading-lamp text-body">Loading book details...</p>;
  }

  const field = (label: string, value: string | undefined) =>
    value ? (
      <div>
        <dt class="text-data text-reading-lamp mb-1">{label}</dt>
        <dd class="text-body text-ink">{value}</dd>
      </div>
    ) : null;

  return (
    <div class="space-y-8">
      <header>
        <a href="/catalog" class="text-reading-lamp text-data hover:text-leather transition-colors">
          &larr; Back to Catalog
        </a>
        <h1 class="font-display text-page-title text-ink mt-2">{book.title}</h1>
        <p class="text-reading-lamp text-data mt-1">by {book.author}</p>
      </header>

      <hr class="gold-rule" aria-hidden="true" />

      <dl class="grid grid-cols-1 md:grid-cols-2 gap-6">
        {field("Genre", book.genre)}
        {field("ISBN", book.isbn)}
        {field("Location", book.location)}
        {field("Status", book.status.replace("_", " "))}
        {field("Added", new Date(book.created_at).toLocaleDateString())}
        {field("Updated", new Date(book.updated_at).toLocaleDateString())}
      </dl>

      <div class="pt-4 border-t border-book-cloth">
        <button
          onClick={handleDelete}
          disabled={deleting}
          class="px-4 py-2 rounded border border-red-300 text-red-700 text-body hover:bg-red-50 transition-colors disabled:opacity-50"
        >
          {deleting ? "Removing..." : "Remove from Catalog"}
        </button>
      </div>
    </div>
  );
}
