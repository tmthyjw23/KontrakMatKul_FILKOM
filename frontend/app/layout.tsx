import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Login — Kontrak Matakuliah FILKOM UNKLAB",
  description:
    "Sistem Kontrak Matakuliah Fakultas Ilmu Komputer, Universitas Klabat. Login sebagai Student atau Admin.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body>{children}</body>
    </html>
  );
}
