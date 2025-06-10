"use client";
import { Montserrat } from "next/font/google";
import "./globals.css";
import { cn } from "../utils/utils";
import useGuard from "@/hooks/useGuard";
import useClientSide from "@/hooks/useClientSide";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const montserrat = Montserrat({
  subsets: ['latin'],
  weight: ['400', '600', '700'],
  variable: '--font-montserrat',
  display: 'swap',
})

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
          `${montserrat.className} antialiased`,

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
