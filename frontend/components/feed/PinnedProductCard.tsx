import type { FeedProduct } from "@/types/commerce";

interface PinnedProductCardProps {
    product: FeedProduct;
    onAddToCart: () => void;
    onVisitStore: () => void;
}

export default function PinnedProductCard({ product, onAddToCart, onVisitStore }: PinnedProductCardProps) {
    return (
        <div className="rounded-[28px] border border-white/10 bg-black/65 p-4 text-white shadow-xl shadow-black/30">
            <div className="flex items-center gap-3">
                <img
                    src={product.image_url}
                    alt={product.name}
                    className="h-16 w-16 rounded-3xl object-cover ring-2 ring-white/20"
                />
                <div className="min-w-0 flex-1">
                    <p className="text-sm font-semibold text-white">{product.name}</p>
                    <p className="mt-1 text-xs text-slate-300">{product.merchant_name}</p>
                </div>
                <div className="text-right">
                    <p className="text-sm font-semibold text-white">${product.price.toFixed(2)}</p>
                    {product.discount ? (
                        <p className="text-[11px] text-emerald-300">{product.discount}% off</p>
                    ) : null}
                </div>
            </div>
            <div className="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <button
                    type="button"
                    className="rounded-full bg-white px-4 py-2 text-sm font-semibold text-slate-900 transition hover:bg-slate-100"
                    onClick={onAddToCart}
                >
                    Add to cart
                </button>
                <button
                    type="button"
                    className="rounded-full border border-white/20 bg-white/5 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/10"
                    onClick={onVisitStore}
                >
                    Visit store
                </button>
            </div>
        </div>
    );
}
