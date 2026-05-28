import Link from "next/link";
import MediaGallery from "@/components/product/MediaGallery";
import ProductCard from "@/components/product/ProductCard";
import ProductPinning from "@/components/product/ProductPinning";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import SectionHeader from "@/components/ui/SectionHeader";
import { products, recommendedProducts } from "@/lib/mock/products";

const parsePrice = (price: string) => Number(price.replace("$", ""));

export default async function ProductPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const selectedProduct = products.find((product) => product.id === id) ?? products[0];
    const price = parsePrice(selectedProduct.price);
    const compareAt = price + 12;
    const media = [selectedProduct.image, selectedProduct.image, selectedProduct.image];
    const relatedProducts = recommendedProducts.filter((product) => product.id !== selectedProduct.id).slice(0, 4);

    return (
        <div className="space-y-8 pb-10">
            <div className="grid gap-8 xl:grid-cols-[1.1fr_0.9fr]">
                <Card className="p-5 sm:p-6">
                    <div className="space-y-2">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-[#E91E63]">Featured product</p>
                        <SectionHeader title={selectedProduct.title ?? selectedProduct.name} description={selectedProduct.description} />
                    </div>
                    <div className="mt-4 grid gap-6">
                        <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                            <MediaGallery images={media} />
                        </div>
                        <div className="space-y-4 rounded-[24px] bg-[#FFF8FB] p-5">
                            <div className="flex flex-wrap gap-2">
                                {selectedProduct.badge ? (
                                    <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">{selectedProduct.badge}</span>
                                ) : null}
                                <span className="rounded-full bg-emerald-100 px-3 py-1 text-[11px] font-semibold text-emerald-800">In stock</span>
                            </div>
                            <div className="flex flex-wrap items-end justify-between gap-4">
                                <div>
                                    <p className="text-3xl font-semibold text-zinc-900">${price}</p>
                                    <p className="text-sm text-zinc-400 line-through">${compareAt}</p>
                                </div>
                                <div className="text-right">
                                    <p className="text-sm font-semibold text-zinc-900">{selectedProduct.rating ?? "4.8"} ★</p>
                                    <p className="text-sm text-zinc-500">{selectedProduct.reviewCount ?? "1.2k reviews"}</p>
                                </div>
                            </div>
                            <p className="text-sm leading-7 text-zinc-600">
                                {selectedProduct.shopperNote ?? selectedProduct.description}
                            </p>
                            <div className="grid gap-2 text-sm text-zinc-700 sm:grid-cols-2">
                                <p><span className="font-semibold">Creator:</span> {selectedProduct.creator ?? "Teamart Studio"}</p>
                                <p><span className="font-semibold">Merchant:</span> {selectedProduct.merchant ?? "@teamart"}</p>
                                <p><span className="font-semibold">Available stock:</span> {selectedProduct.inventory ?? "Ready to ship"}</p>
                                <p><span className="font-semibold">Shipping:</span> {selectedProduct.deliveryWindow ?? "2–4 business days"}</p>
                            </div>
                            <div className="flex flex-wrap gap-3">
                                <Button asChild variant="primary" className="flex-1 min-w-[160px]">
                                    <Link href="/cart">Add to cart</Link>
                                </Button>
                                <Button asChild variant="secondary" className="flex-1 min-w-[160px]">
                                    <Link href="/account/wishlist">Save for later</Link>
                                </Button>
                            </div>
                            <div className="grid gap-3 sm:grid-cols-2">
                                {selectedProduct.highlights?.map((item) => (
                                    <div key={item} className="rounded-[20px] bg-white px-4 py-3 text-sm text-zinc-700">
                                        {item}
                                    </div>
                                ))}
                            </div>
                            <ProductPinning />
                        </div>
                    </div>
                </Card>

                <div className="space-y-4">
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Customer love</p>
                        <h2 className="mt-3 text-xl font-semibold text-zinc-900">Why shoppers are adding it</h2>
                        <ul className="mt-4 space-y-3 text-sm text-zinc-600">
                            <li>• Premium materials with creator-first styling</li>
                            <li>• Strong live-room performance and bundle readiness</li>
                            <li>• Fast-moving in social and storefront discovery</li>
                        </ul>
                    </Card>
                    <Card className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Product details</p>
                        <div className="mt-3 space-y-3 text-sm text-zinc-700">
                            <div className="rounded-[20px] bg-[#FFF8FB] p-4">
                                <p className="font-semibold text-zinc-900">Category</p>
                                <p className="mt-1 text-zinc-600">{selectedProduct.category ?? "Lifestyle"}</p>
                            </div>
                            <div className="rounded-[20px] bg-[#FFF8FB] p-4">
                                <p className="font-semibold text-zinc-900">Shipping</p>
                                <p className="mt-1 text-zinc-600">{selectedProduct.deliveryWindow ?? "Available for 2–4 business day delivery across featured markets."}</p>
                            </div>
                            <div className="rounded-[20px] bg-[#FFF8FB] p-4">
                                <p className="font-semibold text-zinc-900">Variants</p>
                                <div className="mt-2 flex flex-wrap gap-2">
                                    {selectedProduct.variants?.map((variant) => (
                                        <span key={variant.label} className="rounded-full bg-white px-3 py-1 text-[11px] font-semibold text-zinc-700">
                                            {variant.label}: {variant.value}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </Card>
                </div>
            </div>

            <section className="space-y-4">
                <SectionHeader
                    title="Related products"
                    description="Complete the cart with creator favorites and adjacent products your shoppers are testing next."
                />
                <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
                    {relatedProducts.map((product) => (
                        <ProductCard key={product.id} product={product} />
                    ))}
                </div>
            </section>
        </div>
    );
}
