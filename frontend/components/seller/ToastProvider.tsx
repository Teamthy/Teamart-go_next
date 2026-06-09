"use client";

import { createContext, useCallback, useContext, useMemo, useState } from "react";
import { AlertTriangle, CheckCircle2, Info, X } from "lucide-react";
import type { ReactNode } from "react";

interface ToastMessage {
    id: string;
    title: string;
    description?: string;
    type: "success" | "error" | "info";
}

interface ToastContextValue {
    toast: (message: Omit<ToastMessage, "id">) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

const iconMap = {
    success: <CheckCircle2 className="h-5 w-5" />,
    error: <AlertTriangle className="h-5 w-5" />,
    info: <Info className="h-5 w-5" />,
};

const toneMap: Record<ToastMessage["type"], string> = {
    success: "border-emerald-200 bg-emerald-50 text-emerald-900",
    error: "border-rose-200 bg-rose-50 text-rose-900",
    info: "border-sky-200 bg-sky-50 text-sky-900",
};

export function ToastProvider({ children }: { children: ReactNode }) {
    const [toasts, setToasts] = useState<ToastMessage[]>([]);

    const toast = useCallback((message: Omit<ToastMessage, "id">) => {
        const id = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
        setToasts((current) => [...current, { id, ...message }]);
        window.setTimeout(() => {
            setToasts((current) => current.filter((item) => item.id !== id));
        }, 3800);
    }, []);

    const value = useMemo(() => ({ toast }), [toast]);

    return (
        <ToastContext.Provider value={value}>
            {children}
            <div className="fixed bottom-5 right-5 z-50 flex w-full max-w-sm flex-col gap-3">
                {toasts.map((item) => (
                    <div
                        key={item.id}
                        className={`rounded-3xl border p-4 shadow-xl shadow-slate-900/10 ${toneMap[item.type]}`}
                    >
                        <div className="flex items-start gap-3">
                            <span className="mt-0.5">{iconMap[item.type]}</span>
                            <div className="flex-1">
                                <p className="font-semibold text-sm">{item.title}</p>
                                {item.description && <p className="text-sm mt-1 leading-5">{item.description}</p>}
                            </div>
                            <button
                                type="button"
                                onClick={() => setToasts((current) => current.filter((toast) => toast.id !== item.id))}
                                className="text-slate-400 hover:text-slate-600"
                            >
                                <X className="h-4 w-4" />
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </ToastContext.Provider>
    );
}

export function useToast() {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error("useToast must be used within ToastProvider");
    }
    return context;
}
