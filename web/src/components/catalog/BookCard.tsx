export interface Book {
  id: string;
  title: string;
  author: string;
  isbn?: string;
  genre?: string;
  status: string;
  location?: string;
  created_at: string;
  updated_at: string;
}

interface BookCardProps {
  book: Book;
}

export default function BookCard({ book }: BookCardProps) {
  return (
    <a
      href={`/catalog/${book.id}`}
      class="block bg-book-cloth/30 rounded-card p-6 border border-book-cloth/50 shadow-index-card hover:border-leather/40 transition-colors no-underline"
    >
      <h3 class="font-body font-semibold text-ink text-section leading-tight mb-2">{book.title}</h3>
      <p class="text-reading-lamp text-data mb-3">{book.author}</p>
      <div class="flex items-center justify-between">
        <span class="text-small text-reading-lamp">{book.genre || "Uncategorized"}</span>
        <span
          class={[
            "text-small px-2 py-0.5 rounded",
            book.status === "available"
              ? "bg-green-900/20 text-green-700"
              : book.status === "checked_out"
              ? "bg-leather/20 text-leather"
              : "bg-red-900/20 text-red-700",
          ].join(" ")}
        >
          {book.status.replace("_", " ")}
        </span>
      </div>
      {book.isbn && <p class="text-small text-reading-lamp/70 mt-2 font-mono">ISBN: {book.isbn}</p>}
    </a>
  );
}
