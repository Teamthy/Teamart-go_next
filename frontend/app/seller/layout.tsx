import type { Metadata } from "next";

export const metadata: Metadata = {
    title: "Seller Center | Teamart",
};

export default function SellerLayout({ children }: { children: React.ReactNode }) {
    return <>{children}</>;
}
