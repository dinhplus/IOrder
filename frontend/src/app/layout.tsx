import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: "IOrder",
  description: "IOrder — Restaurant ordering platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}): React.JSX.Element {
  return (
    <html lang="en">
      <body className="min-h-screen bg-background font-sans antialiased">
        <header className="border-b border-border bg-background">
          <div className="container mx-auto flex items-center justify-between px-6 py-4">
            <Link href="/" className="text-xl font-bold text-foreground">
              IOrder
            </Link>
            <nav className="flex items-center gap-4 text-sm text-muted-foreground">
              <Link href="/" className="hover:text-foreground transition-colors">
                Home
              </Link>
              <Link href="/menu" className="hover:text-foreground transition-colors">
                Menu
              </Link>
              <Link href="/orders" className="hover:text-foreground transition-colors">
                Orders
              </Link>
              <Link href="/kitchen" className="hover:text-foreground transition-colors">
                Kitchen
              </Link>
              <Link href="/floor-plan" className="hover:text-foreground transition-colors">
                Floor Plan
              </Link>
            </nav>
          </div>
        </header>
        <main className="container mx-auto px-6 py-8">{children}</main>
        <footer className="border-t border-border mt-16">
          <div className="container mx-auto px-6 py-4 text-center text-sm text-muted-foreground">
            © {new Date().getFullYear()} IOrder. Restaurant ordering platform.
          </div>
        </footer>
      </body>
    </html>
  );
}
