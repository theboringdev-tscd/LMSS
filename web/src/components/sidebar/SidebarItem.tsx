import { cn } from "../../lib/utils";

interface SidebarItemProps {
  label: string;
  href: string;
  active?: boolean;
  icon?: string;
}

const iconMap: Record<string, string> = {
  dashboard: "\u2302",
  catalog: "\u2601",
  patrons: "\u263A",
  circulation: "\u21C4",
  fines: "\u00A4",
  reservations: "\u2605",
  reports: "\u2261",
  settings: "\u2699",
  login: "\u2192",
};

export default function SidebarItem({ label, href, active, icon }: SidebarItemProps) {
  const glyph = icon ? iconMap[icon] ?? "\u2022" : "\u2022";

  return (
    <li>
      <a
        href={href}
        class={cn(
          "flex items-center gap-3 px-4 py-2 rounded text-data transition-colors",
          "hover:bg-white/5 hover:text-paper",
          active
            ? "bg-white/10 text-paper font-medium"
            : "text-reading-lamp"
        )}
        aria-current={active ? "page" : undefined}
      >
        <span class="w-5 text-center text-gold-leaf font-mono" aria-hidden="true">
          {glyph}
        </span>
        <span>{label}</span>
      </a>
    </li>
  );
}
