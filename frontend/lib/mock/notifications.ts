export interface NotificationItem {
    id: string;
    title: string;
    detail: string;
    tone?: string;
    unread?: boolean;
    time?: string;
}

export interface NotificationGroup {
    title: string;
    items: NotificationItem[];
}

export const accountNotifications: NotificationItem[] = [
    {
        id: "n1",
        title: "Creator Collaboration Hoodie is back in stock",
        detail: "The best-selling hoodie has been restocked for the next drop window.",
        tone: "info",
        unread: true,
        time: "Just now",
    },
    {
        id: "n2",
        title: "Order O-5021 moved to processing",
        detail: "Your order has moved to the fulfillment team and is being prepared.",
        tone: "success",
        unread: false,
        time: "12 min ago",
    },
    {
        id: "n3",
        title: "Livestream reminder for tonight at 8 PM",
        detail: "A creator-hosted shopping session is about to begin.",
        tone: "default",
        unread: true,
        time: "1 hr ago",
    },
    {
        id: "n4",
        title: "Billing update available",
        detail: "Your next invoice is ready and synced to your wallet summary.",
        tone: "warning",
        unread: false,
        time: "Today",
    },
];

export const adminNotices: NotificationItem[] = [
    { id: "a1", title: "Moderation queue review needed", detail: "3 new reports need a response before the next live drop.", tone: "warning", unread: true, time: "15 min ago" },
    { id: "a2", title: "Payout exception flagged", detail: "One merchant payout is delayed and requires confirmation.", tone: "warning", unread: true, time: "35 min ago" },
    { id: "a3", title: "Live room alert", detail: "Creator Q&A room is near capacity and trending on discovery.", tone: "success", unread: false, time: "1 hr ago" },
];

export const inboxQuickActions = [
    "Manage account preferences",
    "Review saved favorites",
    "Open order activity",
];

export const notificationPreferences = [
    "Live stream reminders and order updates",
    "Weekly creator drops and promotions",
    "Personalized recommendations enabled",
];

export const notificationGroups: NotificationGroup[] = [
    {
        title: "Today",
        items: accountNotifications.slice(0, 2),
    },
    {
        title: "Earlier",
        items: accountNotifications.slice(2),
    },
];

export const adminAlerts = [
    "3 new support tickets require review",
    "2 merchant payouts are pending approval",
    "A live room is approaching capacity",
];

export const accountAlertSummary = [
    "2 unread notifications",
    "1 payment reminder",
    "1 live drop is starting soon",
];
