interface StatBlockProps {
  value: number | string;
  label: string;
}

export default function StatBlock({ value, label }: StatBlockProps) {
  return (
    <div className="bg-book-cloth/30 rounded-card p-8 border border-book-cloth/50 shadow-index-card">
      <p className="font-display text-display text-ink">{value}</p>
      <p className="text-reading-lamp text-data uppercase tracking-wider mt-2">{label}</p>
    </div>
  );
}
