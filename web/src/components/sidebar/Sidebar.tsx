import SidebarItem from "./SidebarItem";

const navItems = [
  { label: "Dashboard", href: "/", icon: "dashboard" },
  { label: "Catalog", href: "/catalog", icon: "catalog" },
  { label: "Patrons", href: "/patrons", icon: "patrons" },
  { label: "Circulation", href: "/circulation", icon: "circulation" },
  { label: "Fines", href: "/fines", icon: "fines" },
  { label: "Reservations", href: "/reservations", icon: "reservations" },
  { label: "Reports", href: "/reports", icon: "reports" },
  { label: "Settings", href: "/settings", icon: "settings" },
];

export default function Sidebar() {
  return (
    <aside class="fixed top-0 left-0 h-screen w-sidebar bg-binding flex flex-col z-50">
      <div class="px-4 py-6 border-b border-white/5">
        <a href="/" class="flex items-center gap-3 no-underline">
          <span class="w-8 h-8 rounded bg-gold-leaf flex items-center justify-center text-binding font-display font-semibold text-sm">
            LM
          </span>
          <div class="flex flex-col">
            <span class="text-paper font-display font-semibold text-section leading-tight">
              LMSS
            </span>
            <span class="text-reading-lamp text-small">Table of Contents</span>
          </div>
        </a>
      </div>

      <nav class="flex-1 overflow-y-auto py-4" aria-label="Table of Contents">
        <ul class="space-y-1 px-3">
          {navItems.map((item) => (
            <SidebarItem
              key={item.href}
              label={item.label}
              href={item.href}
              icon={item.icon}
            />
          ))}
        </ul>
      </nav>

      <div class="border-t border-white/5 px-4 py-4">
        <SidebarItem
          label="Login"
          href="/login"
          icon="login"
        />
      </div>
    </aside>
  );
}
