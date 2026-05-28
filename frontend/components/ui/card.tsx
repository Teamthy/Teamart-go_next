import type { ReactNode } from "react";

interface CardProps {
    children: ReactNode;
    className?: string;
}

export default function Card({ children, className = "" }: CardProps) {
    return (
        <section className={`rounded-[28px] border border-zinc-200 bg-white ${className}`}>
            {children}
        </section>
    );
}
