import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import ProductCard from "@/components/product/ProductCard";
import SectionHeader from "@/components/ui/SectionHeader";
import { recommendedProducts } from "@/lib/mock/products";

const studioHighlights = [
    {
        title: "Schedule drops",
        description: "Plan livestream launches, announce limited releases, and keep your audience engaged in one place.",
    },
    {
        title: "Track performance",
        description: "Review product views, engagement, and checkout momentum with a creator-first overview.",
    },
    {
        title: "Manage storefront",
        description: "Refresh your catalog, showcase best sellers, and publish new products without leaving the studio.",
    },
];

export default function CreatorStudioPage() {
    return (
        <div className="mx-auto max-w-7xl space-y-8 px-4 py-10 sm:px-6 lg:px-8">
            <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
                <div className="space-y-4">
                    <Badge tone="info">Creator studio</Badge>
                    <SectionHeader
                        title="Build your next launch from a polished creator workspace"
                        description="Coordinate drops, highlight your strongest products, and keep your audience updated with a streamlined workflow."
                    />
                </div>
                <div className="flex flex-wrap gap-3">
                    <Button asChild variant="secondary">
                        <Link href="/creator/preview">Preview profile</Link>
                    </Button>
                    <Button asChild variant="primary">
                        <Link href="/products">Create new product</Link>
                    </Button>
                </div>
            </div>

            <div className="grid gap-4 md:grid-cols-3">
                {studioHighlights.map((item) => (
                    <Card key={item.title} className="p-5">
                        <p className="text-sm font-semibold text-slate-900">{item.title}</p>
                        <p className="mt-3 text-sm text-slate-600">{item.description}</p>
                    </Card>
                ))}
            </div>

            <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
                <Card className="p-6">
                    <p className="text-sm font-semibold text-slate-900">Studio snapshot</p>
                    <div className="mt-5 space-y-4 text-sm text-slate-700">
                        <div className="rounded-[20px] border border-slate-200 bg-slate-50 px-4 py-4">
                            <p className="font-semibold text-slate-900">Live streams</p>
                            <p className="mt-1">2 upcoming sessions today</p>
                        </div>
                        <div className="rounded-[20px] border border-slate-200 bg-slate-50 px-4 py-4">
                            <p className="font-semibold text-slate-900">Featured products</p>
                            <p className="mt-1">6 items are currently spotlighted</p>
                        </div>
                        <div className="rounded-[20px] border border-slate-200 bg-slate-50 px-4 py-4">
                            <p className="font-semibold text-slate-900">Audience engagement</p>
                            <p className="mt-1">+18% from the last campaign</p>
                        </div>
                    </div>
                </Card>

                <Card className="p-6">
                    <div className="flex items-center justify-between gap-3">
                        <div>
                            <p className="text-sm font-semibold text-slate-900">Featured products</p>
                            <p className="text-sm text-slate-500">A curated selection of your best-performing catalog items.</p>
                        </div>
                        <Button asChild variant="secondary">
                            <Link href="/products">Browse catalog</Link>
                        </Button>
                    </div>
                    <div className="mt-5 grid gap-4 md:grid-cols-2">
                        {recommendedProducts.slice(0, 4).map((product) => (
                            <ProductCard key={product.id} product={product} />
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}
