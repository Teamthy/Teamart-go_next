import Link from "next/link";

type Product = {
    id: string | number;
    name: string;
    description?: string;
    price: string | number;
    compareAt?: string | number;
    image?: string;
    badge?: string;
    merchant?: string;
    likes?: string;
    comments?: string;
};

export default function ProductCard({ product }: { product: Product }) {
    const price = typeof product.price === "number" ? `$${product.price.toFixed(0)}` : product.price;
    const compareAt = typeof product.compareAt === "number" ? `$${product.compareAt.toFixed(0)}` : product.compareAt;
    const description = product.description ?? "Explore this item.";

    return (
        <article className="group overflow-hidden rounded-[28px] border border-zinc-200 bg-white">
            <div className="relative aspect-[4/3] overflow-hidden bg-[#FCE4EC]">
                {product.image ? (
                    <img
                        src={product.image}
                        alt={product.name}
                        className="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                    />
                ) : (
                    <div className="flex h-full items-end justify-end p-4">
                        <span className="rounded-full bg-white/90 px-3 py-1 text-[11px] font-semibold text-zinc-800">
                            {product.badge ?? "Trending"}
                        </span>
                    </div>
                )}
                <div className="absolute left-3 top-3 rounded-full bg-black/70 px-3 py-1 text-[11px] font-semibold text-white">
                    {product.badge ?? "Live drop"}
                </div>
            </div>

            <div className="space-y-4 p-4 sm:p-5">
                <div className="space-y-2">
                    <div className="flex items-start justify-between gap-3">
                        <div>
                            <h3 className="text-[17px] font-semibold tracking-tight text-zinc-900">{product.name}</h3>
                            <p className="mt-2 text-sm leading-6 text-zinc-500">{description}</p>
                        </div>
                        <div className="text-right">
                            <p className="text-lg font-semibold text-zinc-900">{price}</p>
                            {compareAt ? (
                                <p className="text-sm text-zinc-400 line-through">{compareAt}</p>
                            ) : null}
                        </div>
                    </div>
                </div>

                <div className="flex flex-wrap items-center justify-between gap-3 text-[12px] text-zinc-500">
                    <span>{product.merchant ?? "@teamart"}</span>
                    <div className="flex items-center gap-3">
                        {product.likes ? <span>{product.likes} likes</span> : null}
                        {product.comments ? <span>{product.comments} comments</span> : null}
                    </div>
                </div>

                <Link
                    href={`/products/${product.id}`}
                    className="inline-flex w-full items-center justify-center rounded-[24px] bg-[#E91E63] px-4 py-3 text-sm font-semibold text-white"
                >
                    Shop now
                </Link>
            </div>
        </article>
    );
}
