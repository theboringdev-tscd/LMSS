export interface Patron {
  id: string;
  name: string;
  email: string;
  phone?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

interface PatronCardProps {
  patron: Patron;
}

export default function PatronCard({ patron }: PatronCardProps) {
  return (
    <a
      href={`/patrons/${patron.id}`}
      class="block bg-book-cloth/30 rounded-card p-6 border border-book-cloth/50 shadow-index-card hover:border-leather/40 transition-colors no-underline"
    >
      <h3 class="font-body font-semibold text-ink text-section leading-tight mb-1">{patron.name}</h3>
      <p class="text-reading-lamp text-data mb-2">{patron.email}</p>
      {patron.phone && <p class="text-small text-reading-lamp/70 font-mono">{patron.phone}</p>}
      <div class="mt-3">
        <span
          class={[
            "text-small px-2 py-0.5 rounded",
            patron.active ? "bg-green-900/20 text-green-700" : "bg-red-900/20 text-red-700",
          ].join(" ")}
        >
          {patron.active ? "Active" : "Inactive"}
        </span>
      </div>
    </a>
  );
}
