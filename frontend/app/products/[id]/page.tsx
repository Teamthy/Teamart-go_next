"use client";

import { useMemo } from "react";
import MediaGallery from "@/components/product/MediaGallery";
import ProductPinning from "@/components/product/ProductPinning";
import WishlistButton from "@/components/product/WishlistButton";
import SectionHeader from "@/components/ui/SectionHeader";
import { useProduct } from "@/hooks/useProducts";

export default function ProductPage({ params }: { params: { id: string } }) {
    const { product, isLoading, error } = useProduct(params.id);

    if (isLoading) {
        return <p>Loading product...</p>;
    }

    if (error) {
        return <p className="text-red-600">Error loading product: {error}</p>;
    }

    if (!product) {
        return <p>Product not found</p>;
    }

    const media = product.image_url ? [product.image_url] : ["/images/placeholder-product.png"];
    const priceLabel = typeof product.price === "number" ? `$${product.price}` : String(product.price);
    const merchantName = (product as any).merchant_name || "Teamart";
    const stockLabel = typeof product.stock === "number" ? String(product.stock) : "—";
    const description = product.description || "Explore this product and add it to your cart."

    return (
        <div className="space-y-8">
            <div className="grid gap-8 xl:grid-cols-[0.7fr_0.3fr]">
                <section className="space-y-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
                    <SectionHeader title={product.name} description={description} />
                    <div className="grid gap-6 lg:grid-cols-[0.9fr_0.4fr]">
                        <div>
                            <MediaGallery images={media} />
                        </div>
                        <div className="space-y-6 rounded-3xl border border-slate-200 bg-slate-50 p-6">
                            <div className="space-y-3">
                                <p className="text-sm uppercase tracking-[0.2em] text-slate-500">Price</p>
                                <p className="text-4xl font-semibold text-slate-900">{priceLabel}</p>
                            </div>
                            <div className="space-y-4">
                                <p className="text-sm text-slate-600">Creator: {merchantName}</p>
                                <p className="text-sm text-slate-600">Available stock: {stockLabel}</p>
                            </div>
                            <button className="w-full rounded-3xl bg-slate-900 px-5 py-4 text-sm font-semibold text-white transition hover:bg-slate-700">
                                Add to cart
                            </button>
                            <WishlistButton productId={String(product.id)} />
                            <ProductPinning />
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}
