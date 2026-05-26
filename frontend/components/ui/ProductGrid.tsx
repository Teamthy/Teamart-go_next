import ProductCard from "@/components/product/ProductCard";

interface ProductGridProps {
    products: Array<{
        id: string | number;
        name: string;
        price: string | number;
        description?: string;
        image?: string;
        badge?: string;
        merchant?: string;
        likes?: string;
        comments?: string;
    }>;
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
