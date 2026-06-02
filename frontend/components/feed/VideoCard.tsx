import type { FeedItem } from "@/types/commerce";

interface VideoCardProps {
    item: FeedItem;
    muted: boolean;
    registerVideo: (id: number) => (node: HTMLVideoElement | null) => void;
}

export default function VideoCard({ item, muted, registerVideo }: VideoCardProps) {
    if (item.kind !== "video") {
        return (
            <img
                src={item.image_url ?? item.product.image_url}
                alt={item.title}
                className="h-full w-full object-cover"
            />
        );
    }

    return (
        <video
            ref={registerVideo(item.id)}
            className="h-full w-full object-cover"
            loop
            muted={muted}
            playsInline
            poster={item.image_url ?? item.product.image_url}
            src={item.video_url}
        />
    );
}
