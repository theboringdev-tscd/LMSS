interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  onSearch?: (query: string) => void;
}

export default function SearchBar({ value, onChange, onSearch }: SearchBarProps) {
  return (
    <form
      class="flex gap-3"
      onSubmit={(e) => {
        e.preventDefault();
        onSearch?.(value);
      }}
    >
      <input
        type="text"
        value={value}
        onInput={(e) => onChange((e.target as HTMLInputElement).value)}
        placeholder="Search by title, author, or ISBN..."
        class="flex-1 px-4 py-2 rounded border border-book-cloth bg-paper text-ink text-body focus:outline-none focus:border-leather focus:ring-1 focus:ring-leather"
      />
      <button
        type="submit"
        class="px-4 py-2 rounded bg-leather text-paper font-body font-medium text-body hover:bg-leather/90 transition-colors"
      >
        Search
      </button>
      {value && (
        <button
          type="button"
          onClick={() => {
            onChange("");
            onSearch?.("");
          }}
          class="px-3 py-2 rounded border border-book-cloth text-reading-lamp text-body hover:border-leather transition-colors"
        >
          Clear
        </button>
      )}
    </form>
  );
}
