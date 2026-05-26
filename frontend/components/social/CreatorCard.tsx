import Link from "next/link";

export type CreatorCardProps = {
    name: string;
    handle: string;
    avatar: string;
    followers: string;
    bio: string;
    mutual?: boolean;
    cta?: string;
    href?: string;
};

export default function CreatorCard({
    name,
    handle,
    avatar,
    followers,
    bio,
    mutual = false,
    cta = "Follow",
    href = "/creator/mia",
}: CreatorCardProps) {
    return (
        <article className="rounded-[28px] border border-zinc-200 bg-white p-4">
            <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                    <img src={avatar} alt={name} className="h-12 w-12 rounded-[20px] object-cover" />
                    <div>
                        <p className="text-sm font-semibold text-zinc-900">{name}</p>
                        <p className="text-[12px] text-zinc-500">{handle}</p>
                    </div>
                </div>
                <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">
                    {followers}
                </span>
            </div>
            <p className="mt-3 text-sm leading-6 text-zinc-600">{bio}</p>
            <div className="mt-4 flex items-center justify-between gap-3">
                {mutual ? (
                    <span className="rounded-full bg-zinc-100 px-3 py-1 text-[11px] font-semibold text-zinc-700">
                        Mutual follow
                    </span>
                ) : (
                    <span className="rounded-full bg-[#FFF8FB] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">
                        New creator
                    </span>
                )}
                <Link href={href} className="rounded-[24px] bg-[#E91E63] px-4 py-2 text-sm font-semibold text-white">
                    {cta}
                </Link>
            </div>
        </article>
    );
}
