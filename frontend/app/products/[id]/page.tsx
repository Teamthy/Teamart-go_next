import MediaGallery from "@/components/product/MediaGallery";
import ProductPinning from "@/components/product/ProductPinning";
import SectionHeader from "@/components/ui/SectionHeader";

const product = {
    id: "1",
    name: "Sustainable Canvas Tote",
    price: "$32",
    description: "Everyday carry with recycled canvas and bold prints.",
    image: "/product-tote.jpg",
    media: ["/product-tote.jpg", "/product-tote-2.jpg", "/product-tote-3.jpg"],
};

export default function ProductPage() {
    return (
        <div className="space-y-8">
            <div className="grid gap-8 xl:grid-cols-[0.7fr_0.3fr]">
                <section className="space-y-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
                    <SectionHeader title={product.name} description={product.description} />
                    <div className="grid gap-6 lg:grid-cols-[0.9fr_0.4fr]">
                        <div>
                            <MediaGallery images={product.media} />
                        </div>
                        <div className="space-y-6 rounded-3xl border border-slate-200 bg-slate-50 p-6">
                            <div className="space-y-3">
                                <p className="text-sm uppercase tracking-[0.2em] text-slate-500">Price</p>
                                <p className="text-4xl font-semibold text-slate-900">{product.price}</p>
                            </div>
                            <div className="space-y-4">
                                <p className="text-sm text-slate-600">Creator: Mia Rivera</p>
                                <p className="text-sm text-slate-600">Available stock: 32</p>
                            </div>
                            <button className="w-full rounded-3xl bg-slate-900 px-5 py-4 text-sm font-semibold text-white transition hover:bg-slate-700">
                                Add to cart
                            </button>
                            <ProductPinning />
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}
