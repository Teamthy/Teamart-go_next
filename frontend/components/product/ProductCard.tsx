import Link from "next/link";

type Product = {
    id: string;
    name: string;
    description: string;
    price: string;
    image: string;
};

export default function ProductCard({
    product,
    showDetailLink = true,
}: {
    product: Product;
    showDetailLink?: boolean;
}) {
    return (
        <article className="group overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
            <div className="relative aspect-[4/3] overflow-hidden bg-slate-100">
                <img src={product.image} alt={product.name} className="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
            </div>
            <div className="space-y-3 p-5">
                <div className="flex items-center justify-between">
                    <h3 className="text-lg font-semibold text-slate-900">{product.name}</h3>
                    <span className="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-700">{product.price}</span>
                </div>
                <p className="text-sm leading-6 text-slate-500">{product.description}</p>
                {showDetailLink && (
                    <Link
                        href={`/product/${product.id}`}
                        className="inline-flex items-center rounded-2xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-700"
                    >
                        View details
                    </Link>
                )}
            </div>
        </article>
    );
}
