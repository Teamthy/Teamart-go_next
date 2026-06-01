import SectionHeader from "@/components/ui/SectionHeader";
import { sessionHistory } from "@/lib/mock/users";

export default function SessionsPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Session History" description="Review your past sessions and activities." />
            <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">