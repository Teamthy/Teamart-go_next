import { notFound } from "next/navigation";
import CreatorProfileCard from "@/components/product/CreatorProfileCard";
import ProductCard from "@/components/product/ProductCard";
import { creators, recommendedProducts } from "@/lib/mock-data";

function getCreator(slug: string) {
    return creators.find((creator) => creator.id === slug);
}

export default function CreatorProfilePage({ params }: { params: { slug: string } }) {
    const creator = getCreator(params.slug);

    if (!creator) {
        notFound();
    }

    const featuredProducts = recommendedProducts.slice(0, 6);

    return (
        <div className="min-h-screen space-y-8 bg-[#F9F5F8] px-4 py-10 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl space-y-8">
                <CreatorProfileCard creator={creator} />

                <section className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
                    <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                        <div>
                            <h2 className="text-2xl font-semibold text-slate-900">Featured drops</h2>
                            <p className="text-sm text-slate-600">Shop the creator's latest collection and livestream picks.</p>
                        </div>
                        <span className="rounded-full bg-[#FCE4EC] px-4 py-2 text-sm font-semibold text-[#C2185B]">{creator.liveStatus}</span>
                    </div>
                    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                        {featuredProducts.map((product) => (
                            <ProductCard key={product.id} product={product} />
                        ))}
                    </div>
                </section>
            </div>
        </div>
    );
}
