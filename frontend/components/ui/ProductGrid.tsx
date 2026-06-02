import ProductCard from "@/components/product/ProductCard";
import type { Product } from "@/types/product";

interface ProductGridProps {
    products: Product[];
}

export default function ProductGrid({ products }: ProductGridProps) {
    return (
        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
            {products.map((product) => (
                <ProductCard key={String(product.id)} product={product} />
            ))}
        </div>
    );
}
