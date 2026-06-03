export type FeedPostType = "video" | "image" | "live";

export interface AuthorProfile {
    id: number;
    name: string;
    handle: string;
    avatar_url: string;
    verified?: boolean;
    store_name?: string;
    store_logo_url?: string;
}

export interface FeedProduct {
    id: number;
    name: string;
    image_url: string;
    price: number;
    discount?: number;
    merchant_name: string;
    merchant_logo_url?: string;
    stock?: number;
}

export interface FeedItem {
    id: number;
    kind: FeedPostType;
    title: string;
    caption: string;
    hashtags: string[];
    creator: AuthorProfile;
    product: FeedProduct;
    video_url?: string;
    image_url?: string;
    like_count: number;
    comment_count: number;
    bookmarked: boolean;
    liked: boolean;
    viewers?: number;
    is_live?: boolean;
    live_title?: string;
    live_status?: string;
}

export interface FeedPage {
    items: FeedItem[];
    next_offset?: number;
    has_more?: boolean;
    total?: number;
}

export interface LiveRoomSummary {
    id: string;
    title: string;
    host: AuthorProfile;
    viewers: number;
    status: string;
    badge: string;
    thumbnail_url: string;
    is_live: boolean;
    pinned_product: FeedProduct;
    short_description?: string;
}

export interface LiveChatMessage {
    id: string;
    user_id: number;
    user_name: string;
    avatar_url: string;
    message: string;
    emoji?: string;
    timestamp: string;
}

export interface LiveRoomDetails extends LiveRoomSummary {
    description: string;
    start_time?: string;
    elapsed_seconds?: number;
    host_message?: string;
    reactions?: Record<string, number>;
    messages: LiveChatMessage[];
    pinned_products: FeedProduct[];
}

export interface StoreSummary {
    id: number;
    owner_id: number;
    name: string;
    description: string;
    category: string;
    banner_url: string;
    status: string;
    created_at: string;
    updated_at: string;
    slug: string;
    tagline: string;
    followers: string;
    rating: string;
    live_status: string;
    products: number;
}

export interface StoreDetails extends StoreSummary { }

export interface ProductVariant {
    id: number;
    name: string;
    price: number;
    stock: number;
}

export interface ProductDetail {
    id: number;
    name: string;
    description: string;
    price: number;
    image_url: string;
    category?: string;
    stock: number;
    variants?: ProductVariant[];
    merchant_name?: string;
    merchant_logo_url?: string;
}

export interface CartPayload {
    product_id: number;
    quantity: number;
    variant_id?: number;
}
