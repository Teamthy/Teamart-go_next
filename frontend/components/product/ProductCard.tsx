import Link from "next/link";

type ProductCardProduct = {
    id: string | number;
    name: string;
    description?: string;
    price?: string | number;
    image?: string;
    image_url?: string;
};

function formatPrice(price?: string | number) {
    if (price === undefined || price === null || price === "") {
        return "Price unavailable";
    }

    if (typeof price === "number") {
        return `$${price.toFixed(2)}`;
    }

    return price;
}

export default function ProductCard({ product }: { product: ProductCardProduct }) {
    const imageSrc = product.image ?? product.image_url ?? "";
    const description = product.description ?? "No description provided.";

    return (
        <article className="group overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
            <div className="relative aspect-[4/3] overflow-hidden bg-slate-100">
                <img
                    src={imageSrc}
                    alt={product.name}
                    className="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                />
            </div>
            <div className="space-y-3 p-5">
                <div className="flex items-center justify-between">
                    <h3 className="text-lg font-semibold text-slate-900">{product.name}</h3>
                    <span className="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-700">
                        {formatPrice(product.price)}
                    </span>
                </div>
                <p className="text-sm leading-6 text-slate-500">{description}</p>
                <Link
                    href={`/products/${product.id}`}
                    className="inline-flex items-center rounded-2xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-700"
                >
                    Shop now
                </Link>
            </div>
        </article>
    );
}
