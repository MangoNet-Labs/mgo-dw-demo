"use client";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { cn } from "../utils/utils";
import useGuard from "@/hooks/useGuard";
import useClientSide from "@/hooks/useClientSide";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  useGuard();
  useClientSide();

  return (
    <html lang="en">
      <body
        className={cn(
          `${geistSans.variable} ${geistMono.variable} antialiased`,

          `max-w-[500px] m-auto  h-[100vh] w-full`
        )}
      >
        <div
          className={`bg-[url('/images/bj.png')] bg-cover bg-center h-full w-full`}
        >
          {children}
        </div>

        <ToastContainer />
      </body>
    </html>
  );
}
