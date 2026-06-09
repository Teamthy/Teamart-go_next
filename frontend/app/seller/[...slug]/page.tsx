import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import Card, { CardBody, CardFooter, CardHeader } from "@/components/seller/Card";

interface SellerSectionPageProps {
    params: { slug: string[] };
}

const formatSegment = (segment: string) =>
    segment
        .replace(/-/g, " ")
        .replace(/\b\w/g, (char) => char.toUpperCase());

export default function SellerSectionPage({ params }: SellerSectionPageProps) {
    const sections = params.slug || [];
    const title = sections.map(formatSegment).join(" / ");

    return (
        <SellerAppShell>
            <PageHeader
                title={title || "Seller"}
                description="This section is under active development as part of the seller center expansion."
            />

            <div className="grid gap-6 lg:grid-cols-2">
                <Card>
                    <CardHeader>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">What you can do here</p>
                    </CardHeader>
                    <CardBody>
                        <p className="text-sm text-[var(--text-secondary)]">
                            The seller navigation now resolves gracefully for routes that are not yet fully built. This placeholder keeps the experience intact while the rest of the Seller Center is completed.
                        </p>
                    </CardBody>
                    <CardFooter>
                        <p className="text-sm text-[var(--text-tertiary)]">Use the sidebar to continue exploring seller tools.</p>
                    </CardFooter>
                </Card>

                <Card>
                    <CardHeader>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">Next steps</p>
                    </CardHeader>
                    <CardBody>
                        <ul className="space-y-3 text-sm text-[var(--text-secondary)]">
                            <li>Build dedicated product creation and inventory management pages.</li>
                            <li>Wire up orders filtering and shipment workflows.</li>
                            <li>Add analytics subpages for revenue, creators, and livestream performance.</li>
                        </ul>
                    </CardBody>
                </Card>
            </div>
        </SellerAppShell>
    );
}
