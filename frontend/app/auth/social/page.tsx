"use client";

import { Suspense, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { BASE } from "@/lib/api";

function SocialLoginContent() {
    const searchParams = useSearchParams();

    useEffect(() => {
        const provider = searchParams.get("provider") || "google";
        window.location.assign(`${BASE}/auth/${provider}`);
    }, [searchParams]);

    return (
        <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 text-white">
            <div className="max-w-md rounded-2xl border border-white/10 bg-white/5 px-6 py-8 text-center backdrop-blur">
                <p className="text-sm uppercase tracking-[0.3em] text-white/70">Redirecting</p>
                <h1 className="mt-4 text-2xl font-semibold">Starting your social login</h1>
                <p className="mt-3 text-sm text-white/80">You are being redirected to the configured provider. If nothing happens, check that your backend OAuth route is available.</p>
            </div>
        </div>
    );
}

export default function SocialLoginPage() {
    return (
        <Suspense fallback={null}>
            <SocialLoginContent />
        </Suspense>
    );
}
