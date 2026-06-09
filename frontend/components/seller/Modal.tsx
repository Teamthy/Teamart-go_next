"use client";

import { X } from "lucide-react";
import type { ReactNode } from "react";

interface ModalProps {
    open: boolean;
    title: string;
    description?: string;
    onClose: () => void;
    children: ReactNode;
    footer?: ReactNode;
}

export default function Modal({ open, title, description, onClose, children, footer }: ModalProps) {
    if (!open) {
        return null;
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 px-4 py-6">
            <div className="w-full max-w-2xl overflow-hidden rounded-3xl border border-slate-200/70 bg-[var(--bg-white)] shadow-2xl shadow-slate-900/10">
                <div className="flex items-center justify-between border-b border-[var(--border-light)] px-6 py-5">
                    <div>
                        <h2 className="text-xl font-semibold text-[var(--text-primary)]">{title}</h2>
                        {description && <p className="text-sm text-[var(--text-secondary)] mt-1">{description}</p>}
                    </div>
                    <button
                        type="button"
                        onClick={onClose}
                        className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--border-light)] text-[var(--text-secondary)] transition hover:bg-slate-100"
                        aria-label="Close modal"
                    >
                        <X className="h-4 w-4" />
                    </button>
                </div>
                <div className="p-6">{children}</div>
                {footer && <div className="border-t border-[var(--border-light)] px-6 py-4">{footer}</div>}
            </div>
        </div>
    );
}
