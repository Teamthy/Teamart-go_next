"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import Badge from "@/components/seller/Badge";
import { Play, Calendar, Users, Eye } from "lucide-react";

interface Livestream {
  id: string;
  title: string;
  creatorName: string;
  status: "live" | "scheduled" | "ended";
  viewers: number;
  revenue: number;
  startTime: string;
  duration?: number;
}

const mockLivestreams: Livestream[] = [
  {
    id: "1",
    title: "Product Launch Event",
    creatorName: "Sarah Chen",
    status: "live",
    viewers: 2341,
    revenue: 1245.99,
    startTime: "2024-01-15 19:00",
  },
  {
    id: "2",
    title: "Weekend Flash Sale",
    creatorName: "Mike Johnson",
    status: "scheduled",
    viewers: 0,
    revenue: 0,
    startTime: "2024-01-16 18:00",
  },
  {
    id: "3",
    title: "New Collection Showcase",
    creatorName: "Emma Davis",
    status: "ended",
    viewers: 5600,
    revenue: 3456.78,
    startTime: "2024-01-14 20:00",
    duration: 120,
  },
  {
    id: "4",
    title: "Q&A with Brand Ambassador",
    creatorName: "Alex Rodriguez",
    status: "scheduled",
    viewers: 0,
    revenue: 0,
    startTime: "2024-01-17 19:00",
  },
];

const statusConfig = {
  live: { variant: "danger", label: "🔴 Live" },
  scheduled: { variant: "secondary", label: "Scheduled" },
  ended: { variant: "secondary", label: "Ended" },
} as const;

export default function LivestreamPage() {
  const [selectedTab, setSelectedTab] = useState<"active" | "scheduled" | "ended">(
    "active"
  );

  const filteredStreams = mockLivestreams.filter((stream) => {
    if (selectedTab === "active") return stream.status === "live";
    if (selectedTab === "scheduled") return stream.status === "scheduled";
    if (selectedTab === "ended") return stream.status === "ended";
    return true;
  });

  return (
    <SellerAppShell>
      <PageHeader
        title="Livestream Management"
        description="Manage live shopping events and creator collaborations"
        breadcrumbs={[
          { label: "Dashboard", href: "/seller/dashboard" },
          { label: "Livestream" },
        ]}
        actions={
          <SellerButton variant="primary" size="md" icon={<Play className="w-4 h-4" />}>
            Start Livestream
          </SellerButton>
        }
      />

      {/* Tabs */}
      <div className="mb-6 border-b border-[var(--border-light)]">
        <div className="flex gap-8">
          {[
            { id: "active", label: "Active Streams", icon: "🔴" },
            { id: "scheduled", label: "Scheduled", icon: "📅" },
            { id: "ended", label: "Ended", icon: "✅" },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setSelectedTab(tab.id as any)}
              className={`px-4 py-3 font-medium text-sm border-b-2 transition-colors ${
                selectedTab === tab.id
                  ? "border-[var(--primary)] text-[var(--primary)]"
                  : "border-transparent text-[var(--text-tertiary)] hover:text-[var(--text-secondary)]"
              }`}
            >
              {tab.icon} {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Livestreams Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredStreams.map((stream) => (
          <div
            key={stream.id}
            className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden hover:shadow-lg transition-shadow"
          >
            {/* Thumbnail */}
            <div className="relative bg-gradient-to-br from-[var(--primary)] to-[var(--primary-hover)] h-48 flex items-center justify-center">
              {stream.status === "live" && (
                <div className="absolute top-3 right-3">
                  <Badge variant="danger">🔴 LIVE</Badge>
                </div>
              )}
              <div className="text-center">
                <Play className="w-12 h-12 text-[var(--bg-white)] mx-auto mb-2 opacity-50" />
                <p className="text-[var(--bg-white)] text-sm font-medium opacity-70">
                  {stream.status === "scheduled"
                    ? "Upcoming"
                    : stream.status === "live"
                      ? "Now Live"
                      : "Ended"}
                </p>
              </div>
            </div>

            {/* Content */}
            <div className="p-4">
              <h3 className="font-semibold text-[var(--text-primary)] mb-2 line-clamp-2">
                {stream.title}
              </h3>

              <div className="space-y-3 mb-4">
                <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)]">
                  <Users className="w-4 h-4" />
                  <span>{stream.creatorName}</span>
                </div>

                <div className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2 text-[var(--text-secondary)]">
                    <Eye className="w-4 h-4" />
                    <span>{stream.viewers.toLocaleString()} viewers</span>
                  </div>
                  <span className="font-semibold text-[var(--text-primary)]">
                    ${stream.revenue.toFixed(2)}
                  </span>
                </div>

                <div className="flex items-center gap-2 text-sm text-[var(--text-tertiary)]">
                  <Calendar className="w-4 h-4" />
                  <span>{stream.startTime}</span>
                  {stream.duration && <span>• {stream.duration} min</span>}
                </div>
              </div>

              {/* Actions */}
              <div className="flex gap-2">
                {stream.status === "live" && (
                  <>
                    <SellerButton variant="primary" size="sm" className="flex-1">
                      Manage
                    </SellerButton>
                    <SellerButton variant="danger" size="sm" className="flex-1">
                      End Stream
                    </SellerButton>
                  </>
                )}
                {stream.status === "scheduled" && (
                  <>
                    <SellerButton variant="primary" size="sm" className="flex-1">
                      Start Now
                    </SellerButton>
                    <SellerButton variant="secondary" size="sm" className="flex-1">
                      Edit
                    </SellerButton>
                  </>
                )}
                {stream.status === "ended" && (
                  <SellerButton variant="secondary" size="sm" className="flex-1">
                    View Report
                  </SellerButton>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {filteredStreams.length === 0 && (
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-12 text-center">
          <Play className="w-12 h-12 text-[var(--text-tertiary)] mx-auto mb-4 opacity-50" />
          <p className="text-[var(--text-secondary)] mb-4">
            No {selectedTab} livestreams found
          </p>
          <SellerButton variant="primary">
            Schedule New Livestream
          </SellerButton>
        </div>
      )}
    </SellerAppShell>
  );
}
