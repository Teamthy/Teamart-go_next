import type { ReactNode } from "react";
import TopBar from "./TopBar";
import Sidebar from "./Sidebar";
import { ToastProvider } from "./ToastProvider";

interface SellerAppShellProps {
    children: ReactNode;
}

export default function SellerAppShell({ children }: SellerAppShellProps) {
    return (
        <div className="min-h-screen bg-[var(--bg-light)]">
            <TopBar />
            <Sidebar />

            <ToastProvider>
                {/* Main content area with sidebar offset */}
                <main className="ml-64 mt-16 min-h-[calc(100vh-4rem)] p-6">
                    <div className="max-w-[calc(100%-3rem)] mx-auto">{children}</div>
                </main>
            </ToastProvider>
        </div>
    );
}
