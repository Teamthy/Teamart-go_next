import Link from "next/link";
import Badge from "@/components/ui/badge";
import Card from "@/components/ui/card";

interface StoreCardProps {
    name: string;
    slug: string;
    category: string;
    rating: string;
    banner: string;
    tagline: string;
    live: string;
    products: string;
}

export default function StoreCard({ name, slug, category, rating, banner, tagline, live, products }: StoreCardProps) {
    return (
        <Card className="overflow-hidden">
            <img src={banner} alt={name} className="h-40 w-full object-cover" />
            <div className="p-4 sm:p-5">
                <div className="flex items-start justify-between gap-3">
                    <div>
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{category}</p>
                        <h3 className="mt-2 text-lg font-semibold text-zinc-900">{name}</h3>
                    </div>
                    <Badge tone="success">{rating}</Badge>
                </div>
                <p className="mt-3 text-sm leading-6 text-zinc-600">{tagline}</p>
                <div className="mt-4 flex flex-wrap gap-3 text-sm text-zinc-500">
                    <span>{products} products</span>
                    <span>{live}</span>
                </div>
                <Link href={`/stores/${slug}`} className="mt-4 inline-flex rounded-[24px] bg-[#E91E63] px-4 py-3 text-sm font-semibold text-white">
                    Explore store
                </Link>
            </div>
        </Card>
    );
}
