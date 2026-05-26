export interface LiveRoom {
    id: string;
    name: string;
    time: string;
    status: string;
    viewers?: string;
}

export interface LiveStat {
    label: string;
    value: string;
}

export interface LiveComment {
    author: string;
    message: string;
    tone: string;
}

export const liveStats: LiveStat[] = [
    { label: "Live rooms", value: "12" },
    { label: "Active viewers", value: "1.8k" },
    { label: "Watch time", value: "34m" },
];

export const liveRoomHighlights = [
    "Creator spotlight: Spring drop visuals",
    "Pinned product: Signature Graphic Tee",
    "Community reactions are flowing in",
    "Limited-edition bundle is trending in chat",
];

export const liveRoomAudienceActions = [
    "React with hearts and fire emojis",
    "Ask the creator a question",
    "Open the product tray instantly",
    "Save the featured bundle for later",
];

export const liveStatsSummary: LiveStat[] = [
    { label: "Peak viewers", value: "2.4k" },
    { label: "Engagement", value: "11.6%" },
    { label: "Orders during stream", value: "68" },
];

export const liveTopMoments = [
    "Bundle reveal drove the strongest click-through",
    "Audience questions increased watch time",
    "The featured product was saved in cart by many shoppers",
    "Creator Q&A became the top shared moment",
];

export const liveActionPlan = [
    "Send the recap to the creator team",
    "Push the featured item into curated recommendations",
    "Create a follow-up livestream teaser",
    "Close the loop on post-stream saves",
];

export const liveSchedule: LiveRoom[] = [
    { id: "sunset-drop", name: "Sunset drop", time: "Today • 8 PM", status: "Live soon", viewers: "1.2k expected" },
    { id: "creator-qna", name: "Creator Q&A", time: "Tomorrow • 6 PM", status: "Scheduled", viewers: "850 expected" },
    { id: "weekend-bundle", name: "Weekend bundle reveal", time: "Friday • 7 PM", status: "Scheduled", viewers: "1.6k expected" },
    { id: "limited-restock", name: "Limited edition restock", time: "Saturday • 4 PM", status: "Scheduled", viewers: "900 expected" },
];

export const liveUpNext = [
    { label: "Join the room", href: "/live/room" },
    { label: "Creator livestream", href: "/creator/livestream" },
];

export const liveRoomStatus = [
    { label: "Room health", value: "Stable" },
    { label: "Chat activity", value: "High" },
    { label: "Pinned product", value: "Signature Graphic Tee" },
];

export const liveComments: LiveComment[] = [
    { author: "Nadia", message: "The hoodie looks amazing on camera.", tone: "Positive" },
    { author: "Leo", message: "Can we get a size guide added?", tone: "Question" },
    { author: "Kira", message: "I’m adding the tote to my cart now.", tone: "Purchase intent" },
];

export const livePinnedProducts = [
    "Signature Graphic Tee",
    "Artist Collaboration Hoodie",
    "Live Stream Ring Light",
];
