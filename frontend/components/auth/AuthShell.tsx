import type { ReactNode } from "react";

export default function AuthShell({
    children,
    className = "",
}: {
    children: ReactNode;
    className?: string;
}) {
    return (
        <div className={"min-h-screen bg-[#FCE4EC] px-4 py-6 sm:px-6 sm:py-10 " + className}>
            <div className="mx-auto w-full max-w-[420px] overflow-hidden rounded-[3rem] bg-white shadow-[0_30px_80px_rgba(233,30,99,0.12)]">
                {children}
            </div>
        </div>
    );
}
