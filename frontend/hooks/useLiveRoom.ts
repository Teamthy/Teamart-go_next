"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "@/lib/api";
import type { LiveChatMessage, LiveRoomDetails } from "@/types/commerce";

export function useLiveRoom(roomId: string) {
    const queryClient = useQueryClient();

    const roomQuery = useQuery<LiveRoomDetails, Error>({
        queryKey: ["liveRoom", roomId],
        queryFn: async () => await api.getLiveRoom(roomId),
        enabled: Boolean(roomId),
        staleTime: 1000 * 30,
        refetchOnWindowFocus: false,
    });

    const reactionMutation = useMutation({
        mutationFn: async (reaction: string) => api.reactLiveRoom(roomId, reaction),
        onMutate: async (reaction) => {
            await queryClient.cancelQueries(["liveRoom", roomId]);
            const previousRoom = queryClient.getQueryData<LiveRoomDetails>(["liveRoom", roomId]);
            if (previousRoom) {
                const nextReactions = {
                    ...previousRoom.reactions,
                    [reaction]: (previousRoom.reactions?.[reaction] ?? 0) + 1,
                };

                queryClient.setQueryData<LiveRoomDetails>(["liveRoom", roomId], {
                    ...previousRoom,
                    reactions: nextReactions,
                });
            }
            return { previousRoom };
        },
        onError: (_error, _reaction, context) => {
            if (context?.previousRoom) {
                queryClient.setQueryData(["liveRoom", roomId], context.previousRoom);
            }
        },
        onSettled: () => queryClient.invalidateQueries(["liveRoom", roomId]),
    });

    const chatMutation = useMutation({
        mutationFn: async (message: string) => await api.chatLiveRoom(roomId, message),
        onMutate: async (message) => {
            await queryClient.cancelQueries(["liveRoom", roomId]);
            const previousRoom = queryClient.getQueryData<LiveRoomDetails>(["liveRoom", roomId]);
            if (previousRoom) {
                const newMessage: LiveChatMessage = {
                    id: `optimistic-${Date.now()}`,
                    user_id: 0,
                    user_name: "You",
                    avatar_url: "/avatar-placeholder.png",
                    message,
                    timestamp: new Date().toISOString(),
                };
                queryClient.setQueryData<LiveRoomDetails>(["liveRoom", roomId], {
                    ...previousRoom,
                    messages: [...previousRoom.messages, newMessage],
                });
            }
            return { previousRoom };
        },
        onError: (_error, _message, context) => {
            if (context?.previousRoom) {
                queryClient.setQueryData(["liveRoom", roomId], context.previousRoom);
            }
        },
        onSettled: () => queryClient.invalidateQueries(["liveRoom", roomId]),
    });

    return {
        room: roomQuery.data ?? null,
        isLoading: roomQuery.isLoading,
        isError: roomQuery.isError,
        error: roomQuery.error?.message ?? null,
        sendReaction: reactionMutation.mutate,
        isReacting: reactionMutation.isLoading,
        sendChat: chatMutation.mutate,
        isChatting: chatMutation.isLoading,
        refetchRoom: roomQuery.refetch,
    };
}
