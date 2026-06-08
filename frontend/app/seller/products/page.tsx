"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import DataTable from "@/components/seller/DataTable";
import Badge from "@/components/seller/Badge";
import { Plus, Search } from "lucide-react";

interface Product {
  id: string;
  name: string;
  sku: string;
  category: string;
  price: number;
  inventory: number;
  orders: number;
  status: "active" | "inactive" | "out_of_stock";
}

const mockProducts: Product[] = [
  {
    id: "1",
    name: "Wireless Headphones Pro",
    sku: "WHP-001",
    category: "Electronics",
    price: 199.99,
    inventory: 45,
    orders: 234,
    status: "active",
  },
  {
    id: "2",
    name: "Smart Watch Elite",
    sku: "SWE-002",
    category: "Wearables",
    price: 349.99,
    inventory: 12,
    orders: 156,
    status: "active",
  },
  {
    id: "3",
    name: "Portable Speaker Mini",
    sku: "PSM-003",
    category: "Audio",
    price: 89.99,
    inventory: 0,
    orders: 412,
    status: "out_of_stock",
  },
  {
    id: "4",
    name: "USB-C Hub Universal",
    sku: "UCH-004",
    category: "Accessories",
    price: 49.99,
    inventory: 156,
    orders: 89,
    status: "active",
  },
  {
    id: "5",
    name: "Phone Stand Adjustable",
    sku: "PSA-005",
    category: "Accessories",
    price: 29.99,
    inventory: 0,
    orders: 234,
    status: "out_of_stock",
  },
];

const statusConfig = {
  active: { variant: "success", label: "Active" },
  inactive: { variant: "secondary", label: "Inactive" },
  out_of_stock: { variant: "danger", label: "Out of Stock" },
} as const;

export default function ProductsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedIds, setSelectedIds] = useState<(string | number)[]>([]);
  const [sortColumn, setSortColumn] = useState<string>("");
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");

  const handleSort = (column: string, direction: "asc" | "desc") => {
    setSortColumn(column);
    setSortDirection(direction);
  };

  const filteredProducts = mockProducts.filter(
    (product) =>
      product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      product.sku.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <SellerAppShell>
      <PageHeader
        title="Products"
        description="Manage your product catalog and inventory"
        breadcrumbs={[
          { label: "Dashboard", href: "/seller/dashboard" },
          { label: "Products" },
        ]}
        actions={
          <div className="flex gap-3">
            <SellerButton variant="secondary" size="md">
              Export CSV
            </SellerButton>
            <SellerButton variant="primary" size="md" icon={<Plus className="w-4 h-4" />}>
              Add Product
            </SellerButton>
          </div>
        }
      />

      {/* Filters and Search */}
      <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
          {/* Search */}
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
              Search Products
            </label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-tertiary)]" />
              <input
                type="text"
                placeholder="Search by name or SKU..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
              />
            </div>
          </div>

          {/* Category Filter */}
          <div>
            <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
              Category
            </label>
            <select className="w-full px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]">
              <option value="">All Categories</option>
              <option value="electronics">Electronics</option>
              <option value="wearables">Wearables</option>
              <option value="audio">Audio</option>
              <option value="accessories">Accessories</option>
            </select>
          </div>

          {/* Status Filter */}
          <div>
            <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
              Status
            </label>
            <select className="w-full px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]">
              <option value="">All Status</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="out_of_stock">Out of Stock</option>
            </select>
          </div>
        </div>

        {/* Results Info */}
        {selectedIds.length > 0 && (
          <div className="flex items-center justify-between pt-4 border-t border-[var(--border-light)]">
            <span className="text-sm text-[var(--text-secondary)]">
              {selectedIds.length} product(s) selected
            </span>
            <div className="flex gap-2">
              <SellerButton variant="secondary" size="sm">
                Edit
              </SellerButton>
              <SellerButton variant="danger" size="sm">
                Delete
              </SellerButton>
            </div>
          </div>
        )}
      </div>

      {/* Products Table */}
      <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden">
        <div className="p-6 border-b border-[var(--border-light)]">
          <h3 className="text-lg font-semibold text-[var(--text-primary)]">
            All Products ({filteredProducts.length})
          </h3>
        </div>
        <DataTable
          columns={[
            {
              key: "name",
              label: "Product",
              sortable: true,
              render: (value, row: Product) => (
                <div>
                  <p className="font-medium text-[var(--text-primary)]">
                    {value}
                  </p>
                  <p className="text-xs text-[var(--text-tertiary)]">
                    SKU: {row.sku}
                  </p>
                </div>
              ),
            },
            {
              key: "category",
              label: "Category",
              width: "140px",
              sortable: true,
            },
            {
              key: "price",
              label: "Price",
              width: "100px",
              render: (value) => `$${value.toFixed(2)}`,
              sortable: true,
            },
            {
              key: "inventory",
              label: "Inventory",
              width: "100px",
              render: (value, row: Product) => (
                <span
                  className={
                    row.inventory === 0
                      ? "text-[var(--danger)]"
                      : "text-[var(--text-primary)]"
                  }
                >
                  {value} units
                </span>
              ),
              sortable: true,
            },
            {
              key: "orders",
              label: "Orders",
              width: "100px",
              sortable: true,
            },
            {
              key: "status",
              label: "Status",
              width: "140px",
              render: (value) => (
                <Badge variant={statusConfig[value].variant as any}>
                  {statusConfig[value].label}
                </Badge>
              ),
            },
            {
              key: "id",
              label: "Action",
              width: "120px",
              render: () => (
                <SellerButton variant="ghost" size="sm">
                  Edit
                </SellerButton>
              ),
            },
          ]}
          data={filteredProducts}
          selectable={true}
          selectedIds={selectedIds}
          onSelectChange={setSelectedIds}
          sortColumn={sortColumn}
          sortDirection={sortDirection}
          onSort={handleSort}
          emptyMessage="No products found. Try adjusting your search filters."
        />
      </div>
    </SellerAppShell>
  );
}
