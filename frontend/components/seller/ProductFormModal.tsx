"use client";

import { useMemo, useState, type ChangeEvent } from "react";
import Modal from "./Modal";
import SellerButton from "./Button";
import { useToast } from "./ToastProvider";

interface ProductFormModalProps {
    open: boolean;
    onClose: () => void;
    onCreate: (product: {
        name: string;
        sku: string;
        category: string;
        price: number;
        inventory: number;
    }) => void;
}

const emptyForm = {
    name: "",
    sku: "",
    category: "",
    price: "",
    inventory: "",
};

export default function ProductFormModal({ open, onClose, onCreate }: ProductFormModalProps) {
    const { toast } = useToast();
    const [form, setForm] = useState({ ...emptyForm });
    const [errors, setErrors] = useState<Record<string, string>>({});

    const isValid = useMemo(() => {
        return form.name.trim() !== "" && form.sku.trim() !== "" && form.category.trim() !== "" && Number(form.price) > 0;
    }, [form]);

    const handleChange = (field: string, value: string | number) => {
        setForm((prev) => ({ ...prev, [field]: String(value) }));
        setErrors((prev) => ({ ...prev, [field]: "" }));
    };

    const handleSubmit = () => {
        const nextErrors: Record<string, string> = {};

        if (!form.name.trim()) nextErrors.name = "Product name is required.";
        if (!form.sku.trim()) nextErrors.sku = "SKU is required.";
        if (!form.category.trim()) nextErrors.category = "Category is required.";
        if (!form.price || Number(form.price) <= 0) nextErrors.price = "Enter a valid price.";
        if (!form.inventory || Number(form.inventory) < 0) nextErrors.inventory = "Enter inventory quantity.";

        if (Object.keys(nextErrors).length > 0) {
            setErrors(nextErrors);
            return;
        }

        onCreate({
            name: form.name.trim(),
            sku: form.sku.trim(),
            category: form.category.trim(),
            price: Number(form.price),
            inventory: Number(form.inventory),
        });
        setForm({ ...emptyForm });
        onClose();
        toast({ title: "Product added", description: `${form.name} has been added to your catalog.`, type: "success" });
    };

    return (
        <Modal
            open={open}
            title="Add New Product"
            description="Create a new product listing and sync inventory instantly."
            onClose={onClose}
            footer={
                <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
                    <SellerButton variant="secondary" onClick={onClose}>Cancel</SellerButton>
                    <SellerButton variant="primary" onClick={handleSubmit} disabled={!isValid}>Create Product</SellerButton>
                </div>
            }
        >
            <div className="grid gap-4 sm:grid-cols-2">
                <div>
                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                        Product Name
                    </label>
                    <input
                        value={form.name}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => handleChange("name", event.target.value)}
                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                    />
                    {errors.name && <p className="text-xs text-[var(--danger)] mt-1">{errors.name}</p>}
                </div>
                <div>
                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                        SKU
                    </label>
                    <input
                        value={form.sku}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => handleChange("sku", event.target.value)}
                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                    />
                    {errors.sku && <p className="text-xs text-[var(--danger)] mt-1">{errors.sku}</p>}
                </div>
                <div>
                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                        Category
                    </label>
                    <input
                        value={form.category}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => handleChange("category", event.target.value)}
                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                    />
                    {errors.category && <p className="text-xs text-[var(--danger)] mt-1">{errors.category}</p>}
                </div>
                <div>
                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                        Price
                    </label>
                    <input
                        type="number"
                        value={form.price}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => handleChange("price", event.target.value)}
                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                    />
                    {errors.price && <p className="text-xs text-[var(--danger)] mt-1">{errors.price}</p>}
                </div>
                <div>
                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                        Inventory
                    </label>
                    <input
                        type="number"
                        value={form.inventory}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => handleChange("inventory", event.target.value)}
                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                    />
                    {errors.inventory && <p className="text-xs text-[var(--danger)] mt-1">{errors.inventory}</p>}
                </div>
            </div>
        </Modal>
    );
}
