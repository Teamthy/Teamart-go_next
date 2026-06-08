import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import {
  BarChart3,
  Box,
  Clock,
  Gauge,
  Grid3x3,
  Heart,
  Layout,
  MessageSquare,
  Settings,
  ShoppingCart,
  TrendingUp,
  Users,
  Video,
  Zap,
} from "lucide-react";

interface NavItem {
  href: string;
  label: string;
  icon: React.ReactNode;
  subItems?: { href: string; label: string }[];
}

const navItems: NavItem[] = [
  {
    href: "/seller/dashboard",
    label: "Dashboard",
    icon: <Gauge className="w-5 h-5" />,
  },
  {
    href: "/seller/products",
    label: "Products",
    icon: <Box className="w-5 h-5" />,
    subItems: [
      { href: "/seller/products", label: "All Products" },
      { href: "/seller/products/new", label: "Add Product" },
      { href: "/seller/categories", label: "Categories" },
      { href: "/seller/inventory", label: "Inventory" },
    ],
  },
  {
    href: "/seller/orders",
    label: "Orders",
    icon: <ShoppingCart className="w-5 h-5" />,
    subItems: [
      { href: "/seller/orders", label: "All Orders" },
      { href: "/seller/orders/pending", label: "Pending" },
      { href: "/seller/orders/shipped", label: "Shipped" },
      { href: "/seller/orders/delivered", label: "Delivered" },
    ],
  },
  {
    href: "/seller/livestream",
    label: "Livestream",
    icon: <Video className="w-5 h-5" />,
    subItems: [
      { href: "/seller/livestream/active", label: "Active Streams" },
      { href: "/seller/livestream/scheduled", label: "Scheduled" },
      { href: "/seller/livestream/history", label: "History" },
    ],
  },
  {
    href: "/seller/creators",
    label: "Creators",
    icon: <Users className="w-5 h-5" />,
    subItems: [
      { href: "/seller/creators", label: "All Creators" },
      { href: "/seller/creators/invites", label: "Invitations" },
      { href: "/seller/creators/affiliates", label: "Affiliates" },
    ],
  },
  {
    href: "/seller/campaigns",
    label: "Campaigns",
    icon: <Zap className="w-5 h-5" />,
  },
  {
    href: "/seller/marketing",
    label: "Marketing",
    icon: <TrendingUp className="w-5 h-5" />,
    subItems: [
      { href: "/seller/marketing/promotions", label: "Promotions" },
      { href: "/seller/marketing/coupons", label: "Coupons" },
      { href: "/seller/marketing/ads", label: "Ads" },
    ],
  },
  {
    href: "/seller/analytics",
    label: "Analytics",
    icon: <BarChart3 className="w-5 h-5" />,
    subItems: [
      { href: "/seller/analytics/revenue", label: "Revenue" },
      { href: "/seller/analytics/products", label: "Products" },
      { href: "/seller/analytics/creators", label: "Creators" },
      { href: "/seller/analytics/livestream", label: "Livestream" },
    ],
  },
  {
    href: "/seller/finance",
    label: "Finance",
    icon: <Heart className="w-5 h-5" />,
    subItems: [
      { href: "/seller/finance/payouts", label: "Payouts" },
      { href: "/seller/finance/balance", label: "Balance" },
      { href: "/seller/finance/transactions", label: "Transactions" },
    ],
  },
  {
    href: "/seller/messages",
    label: "Messages",
    icon: <MessageSquare className="w-5 h-5" />,
  },
  {
    href: "/seller/settings",
    label: "Settings",
    icon: <Settings className="w-5 h-5" />,
    subItems: [
      { href: "/seller/settings/store", label: "Store" },
      { href: "/seller/settings/business", label: "Business" },
      { href: "/seller/settings/shipping", label: "Shipping" },
      { href: "/seller/settings/payments", label: "Payments" },
      { href: "/seller/settings/team", label: "Team" },
    ],
  },
];

export default function Sidebar() {
  const pathname = usePathname();
  const [expandedItems, setExpandedItems] = useState<string[]>([]);

  const toggleExpand = (href: string) => {
    setExpandedItems((prev) =>
      prev.includes(href) ? prev.filter((h) => h !== href) : [...prev, href]
    );
  };

  const isActive = (href: string) => {
    return pathname === href || pathname.startsWith(href + "/");
  };

  return (
    <aside className="fixed left-0 top-16 h-[calc(100vh-4rem)] w-64 overflow-y-auto border-r border-[var(--border-light)] bg-[var(--bg-white)]">
      <nav className="space-y-1 px-3 py-4">
        {navItems.map((item) => (
          <div key={item.href}>
            {item.subItems ? (
              <>
                <button
                  onClick={() => toggleExpand(item.href)}
                  className={`w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    isActive(item.href)
                      ? "bg-[var(--primary-lighter)] text-[var(--primary)]"
                      : "text-[var(--text-secondary)] hover:bg-[var(--bg-lighter)]"
                  }`}
                >
                  <span className="flex items-center gap-3">
                    {item.icon}
                    {item.label}
                  </span>
                  <span
                    className={`transition-transform ${
                      expandedItems.includes(item.href) ? "rotate-180" : ""
                    }`}
                  >
                    ▼
                  </span>
                </button>
                {expandedItems.includes(item.href) && (
                  <div className="pl-6 space-y-1 mt-1">
                    {item.subItems.map((subItem) => (
                      <Link
                        key={subItem.href}
                        href={subItem.href}
                        className={`block px-3 py-2 rounded-lg text-sm transition-colors ${
                          isActive(subItem.href)
                            ? "bg-[var(--primary-lighter)] text-[var(--primary)] font-medium"
                            : "text-[var(--text-tertiary)] hover:bg-[var(--bg-lighter)]"
                        }`}
                      >
                        {subItem.label}
                      </Link>
                    ))}
                  </div>
                )}
              </>
            ) : (
              <Link
                href={item.href}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive(item.href)
                    ? "bg-[var(--primary-lighter)] text-[var(--primary)]"
                    : "text-[var(--text-secondary)] hover:bg-[var(--bg-lighter)]"
                }`}
              >
                {item.icon}
                {item.label}
              </Link>
            )}
          </div>
        ))}
      </nav>
    </aside>
  );
}
