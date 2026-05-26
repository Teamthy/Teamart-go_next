import SectionHeader from "@/components/ui/SectionHeader";
import { sessionHistory } from "@/lib/mock/users";

export default function SessionsPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Active sessions" description="Review your signed-in devices and session activity for account safety." />
            <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                <div className="space-y-4">
                    {sessionHistory.map((session) => (
                        <div key={session.id} className="flex flex-col gap-2 rounded-3xl border border-slate-100 bg-slate-50 p-4 sm:flex-row sm:items-center sm:justify-between">
                            <div>
                                <p className="font-medium text-slate-900">{session.device}</p>
                                <p className="text-sm text-slate-600">{session.location}</p>
                            </div>
                            <span className={`rounded-full px-3 py-1 text-sm ${session.active ? "bg-emerald-100 text-emerald-800" : "bg-slate-100 text-slate-700"}`}>
                                {session.active ? "Active" : "Signed out"}
                            </span>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
