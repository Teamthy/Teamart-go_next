import { Bookmark, Heart, MessageCircle, ShoppingBag, Share2 } from "lucide-react";

interface FeedActionsProps {
    liked: boolean;
    bookmarked: boolean;
    likeCount: number;
    commentCount: number;
    onLike: () => void;
    onBookmark: () => void;
    onShare: () => void;
    onComment: () => void;
    onAddToCart: () => void;
}

export default function FeedActions({
    liked,
    bookmarked,
    likeCount,
    commentCount,
    onLike,
    onBookmark,
    onShare,
    onComment,
    onAddToCart,
}: FeedActionsProps) {
    return (
        <div className="flex flex-col items-center gap-4 rounded-full bg-black/45 px-3 py-4 shadow-2xl shadow-black/20">
            <button
                type="button"
                className="inline-flex items-center justify-center rounded-full bg-white/10 p-3 text-white transition hover:bg-white/20"
                aria-label="Like"
                onClick={onLike}
            >
                <Heart className={`h-5 w-5 ${liked ? "text-pink-400" : "text-white"}`} />
                <span className="sr-only">Like</span>
            </button>
            <span className="text-xs font-semibold text-white">{likeCount}</span>
            <button
                type="button"
                className="inline-flex items-center justify-center rounded-full bg-white/10 p-3 text-white transition hover:bg-white/20"
                aria-label="Comments"
                onClick={onComment}
            >
                <MessageCircle className="h-5 w-5 text-white" />
                <span className="sr-only">Comments</span>
            </button>
            <span className="text-xs font-semibold text-white">{commentCount}</span>
            <button
                type="button"
                className="inline-flex items-center justify-center rounded-full bg-white/10 p-3 text-white transition hover:bg-white/20"
                aria-label="Share"
                onClick={onShare}
            >
                <Share2 className="h-5 w-5 text-white" />
                <span className="sr-only">Share</span>
            </button>
            <button
                type="button"
                className="inline-flex items-center justify-center rounded-full bg-white/10 p-3 text-white transition hover:bg-white/20"
                aria-label="Bookmark"
                onClick={onBookmark}
            >
                <Bookmark className={`h-5 w-5 ${bookmarked ? "text-cyan-300" : "text-white"}`} />
                <span className="sr-only">Bookmark</span>
            </button>
            <button
                type="button"
                className="inline-flex items-center justify-center rounded-full bg-white/10 p-3 text-white transition hover:bg-white/20"
                aria-label="Add to cart"
                onClick={onAddToCart}
            >
                <ShoppingBag className="h-5 w-5 text-white" />
                <span className="sr-only">Add to cart</span>
            </button>
        </div>
    );
}
