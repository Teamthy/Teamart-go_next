import type { Creator } from "@/types/creator";
import Link from "next/link";

export default function CreatorProfileCard({
    creator,
}: {
    creator: Creator;
}) {
    return (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-4">
                <img src={creator.avatar} alt={creator.name} className="h-16 w-16 rounded-3xl object-cover" />
                <div>
                    <p className="text-lg font-semibold text-slate-900">{creator.name}</p>
                    <p className="text-sm text-slate-500">{creator.handle}</p>
                </div>
            </div>
            <div className="mt-4 grid gap-3">
                <div className="grid gap-1 rounded-3xl border border-slate-200 bg-slate-50 p-4">
                    <p className="text-sm text-slate-500">Followers</p>
                    <p className="text-xl font-semibold text-slate-900">{creator.followers}</p>
                </div>
                <div className="grid gap-1 rounded-3xl border border-slate-200 bg-slate-50 p-4">
                    <p className="text-sm text-slate-500">Live status</p>
                    <p className="text-xl font-semibold text-slate-900 capitalize">{creator.liveStatus ?? "offline"}</p>
                </div>
            </div>
            <Link
                href={`/creator/${creator.handle}/shop`}
                className="mt-6 inline-flex w-full items-center justify-center rounded-full bg-[#E91E63] px-5 py-3 text-sm font-semibold text-white transition hover:bg-pink-600"
            >
                Visit storefront
            </Link>
        </div>
    );
}
