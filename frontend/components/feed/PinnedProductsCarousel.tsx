import type { FeedProduct } from "@/types/commerce";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PinnedProductsCarouselProps {
    products: FeedProduct[];
    activeIndex: number;
    onSelect: (index: number) => void;
    onPrevious: () => void;
    onNext: () => void;
}

export default function PinnedProductsCarousel({
    products,
    activeIndex,
    onSelect,
    onPrevious,
    onNext,
}: PinnedProductsCarouselProps) {
    return (
        <div className="space-y-4 rounded-3xl border border-slate-200/80 bg-white/90 p-4 shadow-lg shadow-black/5">
            <div className="flex items-center justify-between gap-4">
                <div>
                    <p className="text-xs uppercase tracking-[0.24em] text-slate-500">Pinned products</p>
                    <p className="mt-2 text-sm font-semibold text-slate-900">Shop the room highlights</p>
                </div>
                <div className="flex items-center gap-2">
                    <button
                        type="button"
                        onClick={onPrevious}
                        className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-slate-200 bg-white text-slate-700 transition hover:bg-slate-50"
                        aria-label="Previous product"
                    >
                        <ChevronLeft className="h-4 w-4" />
                    </button>
                    <button
                        type="button"
                        onClick={onNext}
                        className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-slate-200 bg-white text-slate-700 transition hover:bg-slate-50"
                        aria-label="Next product"
                    >
                        <ChevronRight className="h-4 w-4" />
                    </button>
                </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
                {products.map((product, index) => (
                    <button
                        key={product.id}
                        type="button"
                        onClick={() => onSelect(index)}
                        className={`rounded-3xl border px-3 py-3 text-left transition ${index === activeIndex
                                ? "border-pink-500 bg-pink-50"
                                : "border-slate-200 bg-white hover:bg-slate-50"
                            }`}
                    >
                        <p className="text-sm font-semibold text-slate-900">{product.name}</p>
                        <p className="mt-1 text-xs text-slate-500">{product.merchant_name}</p>
                        <p className="mt-2 text-sm font-semibold text-slate-900">${product.price.toFixed(2)}</p>
                    </button>
                ))}
            </div>
        </div>
    );
}
