import Link from "next/link";

type CreatorProfile = {
    id: string;
    name: string;
    handle: string;
    bio: string;
    avatar: string;
    followers: string;
    liveStatus: string;
    products: number;
    rating: number;
    category: string;
};

export default function CreatorProfileCard({ creator }: { creator: CreatorProfile }) {
    return (
        <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex flex-col gap-6 sm:flex-row sm:items-center">
                <img src={creator.avatar} alt={creator.name} className="h-20 w-20 rounded-full object-cover" />
                <div>
                    <p className="text-sm uppercase tracking-[0.32em] text-[#E91E63]">{creator.category}</p>
                    <h3 className="mt-3 text-2xl font-semibold text-slate-900">{creator.name}</h3>
                    <p className="text-sm text-slate-600">@{creator.handle}</p>
                </div>
            </div>

            <p className="mt-6 text-sm leading-6 text-slate-700">{creator.bio}</p>

            <div className="mt-6 grid gap-4 sm:grid-cols-3">
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Followers</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">{creator.followers}</p>
                </div>
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Products</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">{creator.products}</p>
                </div>
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Rating</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">{creator.rating.toFixed(1)}</p>
                </div>
            </div>

            <div className="mt-6 flex flex-wrap items-center gap-3">
                <span className="rounded-full bg-[#FFF0F6] px-3 py-1 text-sm font-semibold text-[#C2185B]">{creator.liveStatus}</span>
                <Link
                    href={`/creator/${creator.id}`}
                    className="inline-flex items-center justify-center rounded-full bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-700"
                >
                    Visit creator page
                </Link>
            </div>
        </section>
    );
}
