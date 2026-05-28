import type { Product } from "@/types/product";
import Link from "next/link";

export default function FeedCard({ product }: { product: Product }) {
    return (
        <article className="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
            <Link href={`/products/${product.id}`} className="block h-[320px] overflow-hidden bg-slate-100">
                <img src={product.image} alt={product.name} className="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
            </Link>
            <div className="space-y-3 p-5">
                <div className="flex items-center justify-between gap-3">
                    <div>
                        <p className="text-sm font-semibold text-slate-900">{product.name}</p>
                        <p className="text-xs text-slate-500">{product.merchant}</p>
                    </div>
                    <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">{product.price.toLocaleString("en-US", { style: "currency", currency: "USD" })}</span>
                </div>
                <p className="text-sm leading-6 text-slate-500">{product.description}</p>
                <Link href={`/products/${product.id}`} className="inline-flex items-center justify-center rounded-full bg-[#E91E63] px-4 py-2 text-sm font-semibold text-white transition hover:bg-pink-600">
                    View shop
                </Link>
            </div>
        </article>
    );
}
