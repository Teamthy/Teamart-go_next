import { renderAccount } from "@/renderers/account";
import { renderAdmin } from "@/renderers/admin";
import { renderAuth } from "@/renderers/auth";
import { renderCreator } from "@/renderers/creator";
import { renderLive } from "@/renderers/live";
import { renderMarketing } from "@/renderers/marketing";
import { renderMerchant } from "@/renderers/merchant";
import { renderProducts } from "@/renderers/products";

export default async function CatchAllPage({ params }: { params: Promise<{ slug?: string[] }> }) {
    const { slug = [] } = await params;
    const route = slug[0] ?? "home";

    switch (route) {
        case "auth":
            return renderAuth(slug);
        case "products":
        case "cart":
        case "checkout":
        case "wishlist":
            return renderProducts(slug);
        case "account":
            return renderAccount(slug);
        case "creator":
            return renderCreator(slug);
        case "live":
            return renderLive(slug);
        case "merchant":
            return renderMerchant(slug);
        case "admin":
            return renderAdmin(slug);
        default:
            return renderMarketing(slug);
    }
}
