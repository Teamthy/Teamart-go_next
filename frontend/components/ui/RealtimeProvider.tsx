"use client";

import { ReactNode } from "react";
import { useRealtime } from "@/hooks/useRealtime";

export default function RealtimeProvider({ children }: { children: ReactNode }) {
    useRealtime();
    return <>{children}</>;
}
