import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Spotify Genre Organizer",
  description: "Organize your Spotify liked songs into genre playlists automatically",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-bg-dark text-text-cream antialiased">
        <div className="grain-overlay" />
        {children}
      </body>
    </html>
  );
}
