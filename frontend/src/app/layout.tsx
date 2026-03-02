import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "IOrder",
  description: "IOrder — Restaurant ordering platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-background font-sans antialiased">
        <header className="border-b px-6 py-4">
          <h1 className="text-xl font-bold">IOrder</h1>
        </header>
        <main className="container mx-auto px-6 py-8">{children}</main>
      </body>
    </html>
  );
}
