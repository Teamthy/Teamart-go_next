"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import DataTable from "@/components/seller/DataTable";
import Badge from "@/components/seller/Badge";
import { Plus, UserPlus } from "lucide-react";

interface Creator {
  id: string;
  name: string;
  username: string;
  followers: number;
  commission: number;
  totalRevenue: number;
  status: "active" | "invited" | "inactive";
  joinedDate: string;
}

const mockCreators: Creator[] = [
  {
    id: "1",
    name: "Sarah Chen",
    username: "@sarahchen",
    followers: 125400,
    commission: 15,
    totalRevenue: 12450.99,
    status: "active",
    joinedDate: "2024-01-01",
  },
  {
    id: "2",
    name: "Mike Johnson",
    username: "@mikej",
    followers: 89300,
    commission: 12,
    totalRevenue: 8234.50,
    status: "active",
    joinedDate: "2024-01-05",
  },
  {
    id: "3",
    name: "Emma Davis",
    username: "@emmadavis",
    followers: 203400,
    commission: 18,
    totalRevenue: 15678.99,
    status: "active",
    joinedDate: "2023-12-15",
  },
  {
    id: "4",
    name: "Alex Rodriguez",
    username: "@alexrod",
    followers: 54200,
    commission: 10,
    totalRevenue: 3450.00,
    status: "invited",
    joinedDate: "2024-01-10",
  },
  {
    id: "5",
    name: "Lisa Wang",
    username: "@lisawang",
    followers: 87600,
    commission: 15,
    totalRevenue: 0,
    status: "inactive",
    joinedDate: "2024-01-08",
  },
];

const statusConfig = {
  active: { variant: "success", label: "Active" },
  invited: { variant: "warning", label: "Invited" },
  inactive: { variant: "secondary", label: "Inactive" },
} as const;

export default function CreatorsPage() {
  const [selectedIds, setSelectedIds] = useState<(string | number)[]>([]);
  const [sortColumn, setSortColumn] = useState<string>("");
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");

  const handleSort = (column: string, direction: "asc" | "desc") => {
    setSortColumn(column);
    setSortDirection(direction);
  };

  const totalRevenue = mockCreators.reduce(
    (sum, creator) => sum + creator.totalRevenue,
    0
  );

  return (
    <SellerAppShell>
      <PageHeader
        title="Creators"
        description="Manage creator partnerships and collaborations"
        breadcrumbs={[
          { label: "Dashboard", href: "/seller/dashboard" },
          { label: "Creators" },
        ]}
        actions={
          <div className="flex gap-3">
            <SellerButton
              variant="secondary"
              size="md"
              icon={<UserPlus className="w-4 h-4" />}
            >
              View Invitations
            </SellerButton>
            <SellerButton variant="primary" size="md" icon={<Plus className="w-4 h-4" />}>
              Invite Creator
            </SellerButton>
          </div>
        }
      />

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <p className="text-sm text-[var(--text-tertiary)] mb-1">
            Total Creators
          </p>
          <p className="text-3xl font-bold text-[var(--text-primary)]">
            {mockCreators.length}
          </p>
        </div>
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <p className="text-sm text-[var(--text-tertiary)] mb-1">
            Active Creators
          </p>
          <p className="text-3xl font-bold text-[var(--success)]">
            {mockCreators.filter((c) => c.status === "active").length}
          </p>
        </div>
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <p className="text-sm text-[var(--text-tertiary)] mb-1">
            Pending Invitations
          </p>
          <p className="text-3xl font-bold text-[var(--warning)]">
            {mockCreators.filter((c) => c.status === "invited").length}
          </p>
        </div>
        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
          <p className="text-sm text-[var(--text-tertiary)] mb-1">
            Creator Revenue
          </p>
          <p className="text-3xl font-bold text-[var(--primary)]">
            ${(totalRevenue / 1000).toFixed(1)}k
          </p>
        </div>
      </div>

      {/* Creators Table */}
      <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden">
        <div className="p-6 border-b border-[var(--border-light)]">
          <h3 className="text-lg font-semibold text-[var(--text-primary)]">
            All Creators
          </h3>
        </div>
        <DataTable
          columns={[
            {
              key: "name",
              label: "Creator",
              sortable: true,
              render: (value, row: Creator) => (
                <div>
                  <p className="font-medium text-[var(--text-primary)]">
                    {value}
                  </p>
                  <p className="text-xs text-[var(--text-tertiary)]">
                    {row.username}
                  </p>
                </div>
              ),
            },
            {
              key: "followers",
              label: "Followers",
              width: "140px",
              render: (value) => `${(value / 1000).toFixed(1)}k`,
              sortable: true,
            },
            {
              key: "commission",
              label: "Commission Rate",
              width: "140px",
              render: (value) => `${value}%`,
              sortable: true,
            },
            {
              key: "totalRevenue",
              label: "Total Revenue",
              width: "140px",
              render: (value) => `$${value.toFixed(2)}`,
              sortable: true,
            },
            {
              key: "status",
              label: "Status",
              width: "120px",
              render: (value) => (
                <Badge variant={statusConfig[value].variant as any}>
                  {statusConfig[value].label}
                </Badge>
              ),
            },
            {
              key: "joinedDate",
              label: "Joined",
              width: "120px",
              sortable: true,
            },
            {
              key: "id",
              label: "Action",
              width: "100px",
              render: () => (
                <SellerButton variant="ghost" size="sm">
                  Manage
                </SellerButton>
              ),
            },
          ]}
          data={mockCreators}
          selectable={true}
          selectedIds={selectedIds}
          onSelectChange={setSelectedIds}
          sortColumn={sortColumn}
          sortDirection={sortDirection}
          onSort={handleSort}
          emptyMessage="No creators found"
        />
      </div>
    </SellerAppShell>
  );
}
