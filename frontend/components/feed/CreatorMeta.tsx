import type { AuthorProfile } from "@/types/commerce";

interface CreatorMetaProps {
    creator: AuthorProfile;
    caption: string;
    hashtags: string[];
    productName: string;
}

export default function CreatorMeta({ creator, caption, hashtags, productName }: CreatorMetaProps) {
    return (
        <div className="space-y-3 rounded-3xl bg-black/55 p-4 text-white shadow-lg shadow-black/10">
            <div className="flex items-center gap-3">
                <img
                    src={creator.avatar_url}
                    alt={creator.name}
                    className="h-12 w-12 rounded-full border border-white/20 object-cover"
                />
                <div>
                    <p className="flex items-center gap-2 text-sm font-semibold text-white">
                        {creator.name}
                        {creator.verified ? (
                            <span className="inline-flex items-center rounded-full bg-white/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.2em] text-white">
                                Verified
                            </span>
                        ) : null}
                    </p>
                    <p className="text-xs text-slate-200">@{creator.handle}</p>
                </div>
            </div>
            <div className="space-y-2 text-sm leading-6 text-slate-100">
                <p>{caption}</p>
                <p className="text-slate-300">Product highlight: <span className="font-semibold text-white">{productName}</span></p>
                <div className="flex flex-wrap gap-2">
                    {hashtags.map((tag) => (
                        <span key={tag} className="rounded-full bg-white/10 px-3 py-1 text-[11px] text-slate-200">
                            {tag}
                        </span>
                    ))}
                </div>
            </div>
        </div>
    );
}
