"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactNode, useState } from "react";
import ErrorBoundary from "./ErrorBoundary";
import RealtimeProvider from "./RealtimeProvider";

export default function Providers({ children }: { children: ReactNode }) {
    const [queryClient] = useState(
        () =>
            new QueryClient({
                defaultOptions: {
                    queries: {
                        staleTime: 1000 * 60 * 2,
                        retry: 2,
                        refetchOnWindowFocus: false,
                        refetchOnReconnect: false,
                    },
                },
            })
    );

    return (
        <ErrorBoundary>
            <QueryClientProvider client={queryClient}>
                <RealtimeProvider>{children}</RealtimeProvider>
            </QueryClientProvider>
        </ErrorBoundary>
    );
}
