import { useState, useEffect } from "react";
import BookCard from "./BookCard";
import SearchBar from "./SearchBar";
import { searchBooks, listBooks, type Book } from "../../lib/api";

interface BookListProps {
  initialBooks?: Book[];
}

export default function BookList({ initialBooks }: BookListProps) {
  const [books, setBooks] = useState<Book[]>(initialBooks ?? []);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (initialBooks) {
      setBooks(initialBooks);
      return;
    }
    setLoading(true);
    listBooks()
      .then(setBooks)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load catalog"))
      .finally(() => setLoading(false));
  }, [initialBooks]);

  const handleSearch = async (q: string) => {
    setQuery(q);
    if (!q.trim()) {
      listBooks()
        .then(setBooks)
        .catch((err) => setError(err instanceof Error ? err.message : "Search failed"));
      return;
    }
    setLoading(true);
    setError(null);
    searchBooks(q)
      .then(setBooks)
      .catch((err) => setError(err instanceof Error ? err.message : "Search failed"))
      .finally(() => setLoading(false));
  };

  if (error) {
    return <p class="text-red-700 text-body bg-red-50 rounded p-4">{error}</p>;
  }

  return (
    <div class="space-y-6">
      <SearchBar value={query} onChange={setQuery} onSearch={handleSearch} />
      {loading && (
        <div class="flex items-center justify-center py-8">
          <div class="h-6 w-6 animate-spin rounded-full border-2 border-book-cloth border-t-leather" />
          <span class="ml-2 text-reading-lamp text-body">Loading...</span>
        </div>
      )}
      {!loading && books.length === 0 && (
        <p class="text-reading-lamp text-body">No books found. Add your first book to get started.</p>
      )}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {books.map((book) => (
          <BookCard key={book.id} book={book} />
        ))}
      </div>
    </div>
  );
}
