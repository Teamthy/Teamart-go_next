"use client";

import { useQuery } from "@tanstack/react-query";
import * as api from "@/lib/api";
import type { LiveRoomSummary } from "@/types/commerce";

export function useLiveFeed() {
    const query = useQuery<LiveRoomSummary[], Error>({
        queryKey: ["liveRooms"],
        queryFn: async () => await api.getLiveRooms(),
        staleTime: 1000 * 60,
        refetchOnWindowFocus: false,
    });

    return {
        rooms: query.data ?? [],
        isLoading: query.isLoading,
        isError: query.isError,
        error: query.isError ? query.error?.message ?? "Unable to load live rooms" : null,
        refetch: query.refetch,
    };
}
