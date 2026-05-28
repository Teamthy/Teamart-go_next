import Link from "next/link";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";

interface RoleCardProps {
    title: string;
    description: string;
    requirements: string[];
    ctaLabel: string;
    href?: string;
    onClick?: () => void;
    tone?: "default" | "success" | "warning";
}

export default function RoleCard({ title, description, requirements, ctaLabel, href, onClick, tone = "default" }: RoleCardProps) {
    return (
        <Card className="p-5 sm:p-6">
            <div className="flex items-start justify-between gap-3">
                <div>
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Role access</p>
                    <h2 className="mt-3 text-xl font-semibold text-zinc-900">{title}</h2>
                </div>
                <Badge tone={tone === "warning" ? "warning" : tone === "success" ? "success" : "default"}>{ctaLabel}</Badge>
            </div>
            <p className="mt-3 text-sm leading-6 text-zinc-600">{description}</p>
            <ul className="mt-4 space-y-2 text-sm text-zinc-700">
                {requirements.map((item) => (
                    <li key={item}>• {item}</li>
                ))}
            </ul>
            {href ? (
                <Button asChild variant="primary" className="mt-5 w-full">
                    <Link href={href}>{ctaLabel}</Link>
                </Button>
            ) : (
                <Button variant="primary" className="mt-5 w-full" onClick={onClick}>
                    {ctaLabel}
                </Button>
            )}
        </Card>
    );
}
