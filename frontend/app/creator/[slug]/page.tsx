import CreatorProfileCard from "@/components/product/CreatorProfileCard";
import ProductCard from "@/components/product/ProductCard";
import { recommendedProducts } from "@/lib/mock-data";

export default function CreatorProfilePage() {
    return (
        <div className="space-y-8">
            <CreatorProfileCard />
            <section className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
                <div className="mb-6 flex items-center justify-between">
                    <div>
                        <h2 className="text-2xl font-semibold text-slate-900">Creator products</h2>
                        <p className="text-sm text-slate-600">Browse featured drops from this creator.</p>
                    </div>
                </div>
                <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
                    {recommendedProducts.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            </section>
        </div>
    );
}
