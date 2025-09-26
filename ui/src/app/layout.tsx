import type { Metadata } from "next";
import { Geist } from "next/font/google";
import "./globals.css";
import { TooltipProvider } from "@/components/ui/tooltip";
import { AgentsProvider } from "@/components/AgentsProvider";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { ThemeProvider } from "@/components/ThemeProvider";
import { Toaster } from "@/components/ui/sonner";
import { AppInitializer } from "@/components/AppInitializer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "kagent.dev",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${geistSans.className} flex flex-col h-screen overflow-hidden`} suppressHydrationWarning>
        <ThemeProvider 
          attribute="class" 
          defaultTheme="system" 
          enableSystem 
          disableTransitionOnChange
          storageKey="kagent-theme"
        >
          <TooltipProvider>
            <AgentsProvider>
              <AppInitializer>
                <Header />
                <main className="flex-1 overflow-y-scroll w-full mx-auto">{children}</main>
                <Footer />
              </AppInitializer>
              <Toaster richColors/>
            </AgentsProvider>
          </TooltipProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
