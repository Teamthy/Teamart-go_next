"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import Input from "@/components/seller/Input";
import SellerButton from "@/components/seller/Button";
import { useToast } from "@/components/seller/ToastProvider";
import { createProduct } from "@/lib/sellerApi";

export default function NewProductPage() {
    const { toast } = useToast();
    const [isSaving, setIsSaving] = useState(false);
    const [form, setForm] = useState({
        name: "",
        sku: "",
        category: "",
        price: "",
        inventory: "",
        description: "",
    });

    const handleChange = (field: keyof typeof form, value: string) => {
        setForm((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async () => {
        const previous = { ...form };
        // Optimistically reset form and show success toast
        setIsSaving(true);
        setForm({ name: "", sku: "", category: "", price: "", inventory: "", description: "" });
        toast({ title: "Creating product", description: `${previous.name || "New product"} — saving...`, type: "info" });
        try {
            await createProduct({
                name: previous.name,
                sku: previous.sku,
                category: previous.category,
                price: previous.price,
                inventory: previous.inventory,
                description: previous.description,
            });
            toast({ title: "Product created", description: `${previous.name || "New product"} has been added.`, type: "success" });
        } catch (e) {
            // rollback form on error
            setForm(previous);
            toast({ title: "Failed", description: "Could not create product.", type: "error" });
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <SellerAppShell>
            <PageHeader
                title="Add New Product"
                description="Create a product listing and manage inventory from the Seller Center."
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Products", href: "/seller/products" },
                    { label: "Add Product" },
                ]}
                actions={
                    <SellerButton variant="secondary" size="md" onClick={() => setForm({ name: "", sku: "", category: "", price: "", inventory: "", description: "" })}>
                        Reset Form
                    </SellerButton>
                }
            />

            <div className="grid gap-6 lg:grid-cols-[1fr_360px]">
                <div className="space-y-6">
                    <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                        <h2 className="text-xl font-semibold text-[var(--text-primary)] mb-3">Product details</h2>
                        <p className="text-sm text-[var(--text-secondary)]">
                            Provide the product title, SKU, category, price, inventory, and a short description.
                        </p>
                    </div>

                    <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                        <div className="grid gap-4">
                            <Input
                                label="Product Name"
                                value={form.name}
                                onValueChange={(value) => handleChange("name", value)}
                            />
                            <Input
                                label="SKU"
                                value={form.sku}
                                onValueChange={(value) => handleChange("sku", value)}
                            />
                            <Input
                                label="Category"
                                value={form.category}
                                onValueChange={(value) => handleChange("category", value)}
                            />
                            <div className="grid gap-4 sm:grid-cols-2">
                                <Input
                                    label="Price"
                                    type="number"
                                    value={form.price}
                                    onValueChange={(value) => handleChange("price", value)}
                                />
                                <Input
                                    label="Inventory"
                                    type="number"
                                    value={form.inventory}
                                    onValueChange={(value) => handleChange("inventory", value)}
                                />
                            </div>
                            <label className="block text-sm font-medium text-[var(--text-secondary)]">Description</label>
                            <textarea
                                value={form.description}
                                onChange={(event) => handleChange("description", event.target.value)}
                                rows={6}
                                className="w-full rounded-xl border border-[var(--border-light)] bg-[var(--bg-white)] px-4 py-3 text-sm text-[var(--text-primary)] focus:outline-none focus:border-[var(--primary)]"
                            />
                        </div>
                    </div>
                </div>

                <div className="space-y-6">
                    <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                        <h2 className="text-xl font-semibold text-[var(--text-primary)] mb-3">Publish settings</h2>
                        <div className="space-y-4 text-sm text-[var(--text-secondary)]">
                            <p>Once saved, the product will appear in your catalog and be available for livestream and marketplace selling.</p>
                            <p>Use the inventory number to keep stock levels updated automatically.</p>
                        </div>
                    </div>

                    <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                        <h2 className="text-xl font-semibold text-[var(--text-primary)] mb-3">Actions</h2>
                        <div className="space-y-3">
                            <div className="rounded-2xl bg-[var(--bg-lighter)] p-4 text-sm text-[var(--text-secondary)]">
                                Tip: Use a clear product title and matching SKU for easier order reconciliation.
                            </div>
                            <SellerButton variant="primary" size="lg" onClick={handleSubmit}>
                                Save Product
                            </SellerButton>
                        </div>
                    </div>
                </div>
            </div>
        </SellerAppShell>
    );
}
