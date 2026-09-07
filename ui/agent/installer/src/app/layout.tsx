import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Config } from "@config/index";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});
export const metadata: Metadata = {
  title: Config.Apps.Title,
  description: Config.Apps.Description,
};


export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html
      lang="en"
      suppressHydrationWarning={true}
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body >
        {/* <div className="flex flex-col justify-center items-center bg-[var(--bg-base)] text-[var(--text-main)] px-4 py-12"> */}
          {children}
        {/* </div> */}
      </body>
    </html>
  );
}