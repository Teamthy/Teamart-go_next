export interface Product {
    id: string | number;
    name: string;
    description: string;
    price: string | number;
    image?: string;
    merchant?: string;
    badge?: string;
    merchantHandle?: string;
    likes?: string;
    comments?: string;
    title?: string;
}
