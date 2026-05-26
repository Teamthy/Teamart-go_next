export interface AdminMetric {
    label: string;
    value: string;
}

export interface AdminRole {
    label: string;
    description: string;
}

export interface AdminTicket {
    id: string;
    title: string;
    priority: string;
    status: string;
}

export interface AuditEntry {
    id: string;
    action: string;
    actor: string;
    time: string;
}

export interface ComplianceItem {
    title: string;
    detail: string;
    status: string;
}

export const adminMetrics: AdminMetric[] = [
    { label: "Pending reviews", value: "18" },
    { label: "Flagged items", value: "7" },
    { label: "System health", value: "99.9%" },
];

export const adminAnalyticsMetrics: AdminMetric[] = adminMetrics;

export const adminOperationalSnapshot = [
    "Review creator compliance actions",
    "Monitor merchant payout exceptions",
    "Confirm moderation queue health",
];

export const adminShortcuts = [
    { label: "Users", href: "/admin/users" },
    { label: "Analytics", href: "/admin/analytics" },
    { label: "Alerts", href: "/admin/alerts" },
    { label: "Tickets", href: "/admin/tickets" },
];

export const adminUserRoles: AdminRole[] = [
    { label: "Creator operators", description: "Role and access view for the Teamart platform" },
    { label: "Merchant admins", description: "Role and access view for the Teamart platform" },
    { label: "Moderation support", description: "Role and access view for the Teamart platform" },
    { label: "Platform analysts", description: "Role and access view for the Teamart platform" },
];

export const moderationQueue = [
    "Live stream report needs review",
    "Product listing flagged for policy compliance",
    "Creator post requires follow-up",
];

export const moderationEscalation = [
    "Review users",
    "See analytics",
];

export const adminSettingsCards = [
    { label: "Feature flags", value: "6 enabled" },
    { label: "Ops status", value: "Healthy" },
    { label: "Release channel", value: "Stable" },
    { label: "Audit retention", value: "180 days" },
];

export const adminAlerts = [
    "3 new support tickets require review",
    "2 merchant payouts are pending approval",
    "A live room is approaching capacity",
];

export const adminTickets: AdminTicket[] = [
    { id: "T-1021", title: "Shipping delay on order O-5018", priority: "High", status: "Open" },
    { id: "T-1022", title: "Creator payout discrepancy", priority: "Medium", status: "In review" },
    { id: "T-1023", title: "Refund exception on O-6034", priority: "High", status: "Escalated" },
];

export const adminAuditEntries: AuditEntry[] = [
    { id: "A-401", action: "Approved payout release", actor: "Avery Lane", time: "09:12" },
    { id: "A-402", action: "Flagged listing for review", actor: "Sofia Kim", time: "10:05" },
    { id: "A-403", action: "Updated live room promo", actor: "Noah Patel", time: "11:18" },
];

export const adminComplianceItems: ComplianceItem[] = [
    { title: "Creator policy review", detail: "All creator bundles have completed checks", status: "Complete" },
    { title: "Merchant refund policy", detail: "Refund exceptions are within the approved SLA", status: "Monitoring" },
    { title: "Audit retention", detail: "Last review completed within the current monthly window", status: "Healthy" },
];

export const adminReports = [
    "Merchant health report ready",
    "Live commerce retention snapshot complete",
    "Moderation quality summary updated",
];
