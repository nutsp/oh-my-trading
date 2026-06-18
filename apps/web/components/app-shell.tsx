import Link from "next/link";

const navigation = [
  ["Dashboard", "/"],
  ["Markets", "/markets/XAUUSD"],
  ["MT5", "/mt5"],
  ["Signals", "/signals"],
  ["Risk", "/risk"],
  ["Journal", "/journal"],
  ["Backtests", "/backtests"],
  ["Agent", "/agent"],
  ["Settings", "/settings"],
] as const;

export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <strong>Oh My Trading</strong>
          <span>AI agent dashboard</span>
        </div>
        <nav className="nav" aria-label="Main navigation">
          {navigation.map(([label, href]) => (
            <Link key={href} href={href} aria-current={href === "/" ? "page" : undefined}>
              <span>{label}</span>
              <span aria-hidden="true">›</span>
            </Link>
          ))}
        </nav>
      </aside>
      <main className="main">
        <header className="topbar">
          <p>XAUUSD, BTC, Forex, Crypto</p>
          <p>Risk mode: manual approval</p>
        </header>
        <section className="content">{children}</section>
      </main>
    </div>
  );
}
