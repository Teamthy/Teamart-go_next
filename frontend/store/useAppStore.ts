import { create } from "zustand";

export interface NotificationItem {
    id: string;
    title: string;
    body: string;
    type: string;
    payload?: any;
    receivedAt: string;
    read?: boolean;
}

export interface FeedUpdateItem {
    id: number;
    name: string;
    description?: string;
    price: number;
    image_url?: string;
    category?: string;
    score?: number;
}

interface AppState {
    notifications: NotificationItem[];
    liveFeedCount: number;
    socketStatus: "idle" | "connecting" | "open" | "closed" | "error";
    addNotification: (notification: NotificationItem) => void;
    addFeedUpdates: (count: number) => void;
    resetLiveFeedCount: () => void;
    setSocketStatus: (status: AppState["socketStatus"]) => void;
}

const useAppStore = create<AppState>((set) => ({
    notifications: [],
    liveFeedCount: 0,
    socketStatus: "idle",
    addNotification: (notification) =>
        set((state) => ({
            notifications: [notification, ...state.notifications].slice(0, 30),
        })),
    addFeedUpdates: (count) =>
        set((state) => ({
            liveFeedCount: state.liveFeedCount + count,
        })),
    resetLiveFeedCount: () => set({ liveFeedCount: 0 }),
    setSocketStatus: (status) => set({ socketStatus: status }),
}));

export default useAppStore;
