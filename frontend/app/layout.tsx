import type { Metadata } from "next";
import { Asap } from "next/font/google";
import "./globals.css";

const geist = Asap({
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "EMX Databas Test",
  description: "EMX Postgres Test",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className={`${geist.className} antialiased w-full min-h-dvh flex flex-col justify-start items-center`}>{children}</body>
    </html>
  );
}
