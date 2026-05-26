# Catch-all Route Documentation

## Summary

- **Distinct page renderers wired in the catch-all dispatcher:** 66
- **Concrete frontend route URLs currently mapped:** 88
- **Source file:** [frontend/app/[...slug]/page.tsx](../frontend/app/[...slug]/page.tsx)

> These routes are **frontend page routes**, not backend API endpoints. Each route is rendered by the catch-all page using shared UI primitives and, where applicable, local mock data.

## Route inventory

### Marketing / static pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Home | `/` | Frontend renderer only | Default route handled by the catch-all page |
| About | `/about` | Frontend renderer only | Marketing page |
| Pricing | `/pricing` | Frontend renderer only | Pricing page |
| Contact | `/contact` | Frontend renderer only | Marketing page |
| FAQ | `/faq` | Frontend renderer only | Marketing page |
| Terms | `/terms` | Frontend renderer only | Marketing page |
| Privacy | `/privacy` | Frontend renderer only | Marketing page |
| Explore | `/explore` | Frontend renderer only | Discovery page |
| Categories | `/categories` | Frontend renderer only | Categories index |
| Category detail | `/categories/[slug]` | Frontend renderer only | Dynamic category view |
| Maintenance | `/maintenance` | Frontend renderer only | Maintenance page |
| Gift cards | `/gift-cards` | Frontend renderer only | Gift card landing |
| Coupons | `/coupon` | Frontend renderer only | Coupon landing |
| Returns | `/returns` | Frontend renderer only | Returns support page |
| Search | `/search` | Frontend renderer only | Search/explore routing |

### Auth / products / commerce

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Auth | `/auth` | Frontend renderer only | Uses auth mode fallback |
| Auth mode | `/auth/[mode]` | Frontend renderer only | Dynamic auth screens |
| Auth social | `/auth/social/[provider]` | Frontend renderer only | Social auth flow |
| Products | `/products` | Frontend renderer only | Catalog landing |
| Create product | `/products/new` | Frontend renderer only | Product builder |
| Drafts | `/products/drafts` | Frontend renderer only | Draft workspace |
| Featured | `/products/featured` | Frontend renderer only | Featured collection |
| Trending | `/products/trending` | Frontend renderer only | Trending collection |
| Recommended | `/products/recommended` | Frontend renderer only | Recommendation rail |
| Recent | `/products/recent` | Frontend renderer only | Recently viewed |
| Deals | `/products/deals` | Frontend renderer only | Promo collection |
| Saved | `/products/saved` | Frontend renderer only | Saved items |
| Compare | `/products/compare` | Frontend renderer only | Product comparison |
| Reviews | `/products/reviews` | Frontend renderer only | Product review view |
| Questions | `/products/questions` | Frontend renderer only | Q&A view |
| Related | `/products/related` | Frontend renderer only | Related products |
| Share | `/products/share` | Frontend renderer only | Share view |
| Edit product | `/products/edit/[id]` | Frontend renderer only | Product editor |
| Cart | `/cart` | Frontend renderer only | Cart summary |
| Cart summary | `/cart/summary` | Frontend renderer only | Summary view |
| Checkout | `/checkout` | Frontend renderer only | Checkout flow entry |
| Checkout shipping | `/checkout/shipping` | Frontend renderer only | Shipping step |
| Checkout payment | `/checkout/payment` | Frontend renderer only | Payment step |
| Checkout review | `/checkout/review` | Frontend renderer only | Review step |
| Checkout success | `/checkout/success` | Frontend renderer only | Success state |
| Checkout failed | `/checkout/failed` | Frontend renderer only | Error state |
| Wishlist | `/wishlist` | Frontend renderer only | Wishlist view |

### Account pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Profile | `/account/profile` | Frontend renderer only | Profile details |
| Addresses | `/account/addresses` | Frontend renderer only | Shipping and billing addresses |
| Billing | `/account/billing` | Frontend renderer only | Billing overview |
| Downloads | `/account/downloads` | Frontend renderer only | Download center |
| Preferences | `/account/preferences` | Frontend renderer only | Preference settings |
| Support | `/account/support` | Frontend renderer only | Support center |
| Security | `/account/security` | Frontend renderer only | Security controls |
| Orders | `/account/orders` | Frontend renderer only | Order history |
| Notifications | `/account/notifications` | Frontend renderer only | Notification center |
| Payments | `/account/payments` | Frontend renderer only | Payment methods |
| Saved items | `/account/saved-items` | Frontend renderer only | Saved products |
| Returns | `/account/returns` | Frontend renderer only | Return requests |
| Subscriptions | `/account/subscriptions` | Frontend renderer only | Subscription management |
| Wallet | `/account/wallet` | Frontend renderer only | Wallet and payout activity |

### Creator pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Creator home | `/creator` | Frontend renderer only | Creator workspace entry |
| Analytics | `/creator/analytics` | Frontend renderer only | Creator analytics |
| Products | `/creator/products` | Frontend renderer only | Creator product manager |
| Livestream | `/creator/livestream` | Frontend renderer only | Livestream hub |
| Studio | `/creator/studio` | Frontend renderer only | Studio landing |
| Studio overview | `/creator/studio/overview` | Frontend renderer only | Studio overview page |
| Studio launches | `/creator/studio/launches` | Frontend renderer only | Studio launch planning |

### Live pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Live home | `/live` | Frontend renderer only | Live shopping hub |
| Live room | `/live/room` | Frontend renderer only | Room-specific live view |

### Merchant pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Merchant home | `/merchant` | Frontend renderer only | Merchant portal entry |
| Analytics | `/merchant/analytics` | Frontend renderer only | Merchant analytics |
| Products | `/merchant/products` | Frontend renderer only | Merchant catalog |
| Inventory | `/merchant/inventory` | Frontend renderer only | Inventory management |
| Fulfillment | `/merchant/fulfillment` | Frontend renderer only | Fulfillment tracker |
| Returns | `/merchant/returns` | Frontend renderer only | Merchant returns |
| Payouts | `/merchant/payouts` | Frontend renderer only | Payouts view |
| Billing | `/merchant/billing` | Frontend renderer only | Billing details |
| Settings | `/merchant/settings` | Frontend renderer only | Merchant settings |
| Customers | `/merchant/customers` | Frontend renderer only | Customer list |
| Teams | `/merchant/teams` | Frontend renderer only | Team visibility |
| Orders | `/merchant/orders` | Frontend renderer only | Merchant orders list |
| Order detail | `/merchant/orders/[id]` | Frontend renderer only | Dynamic order detail |

### Admin pages

| Route | URL | Frontend API surface | Notes |
| --- | --- | --- | --- |
| Admin home | `/admin` | Frontend renderer only | Admin dashboard entry |
| Analytics | `/admin/analytics` | Frontend renderer only | Admin analytics |
| Audit | `/admin/audit` | Frontend renderer only | Audit log |
| Notifications | `/admin/notifications` | Frontend renderer only | Admin notifications |
| Alerts | `/admin/alerts` | Frontend renderer only | Alert center |
| Tickets | `/admin/tickets` | Frontend renderer only | Support tickets |
| Moderation | `/admin/moderation` | Frontend renderer only | Moderation queue |
| Settings | `/admin/settings` | Frontend renderer only | Admin settings |
| Permissions | `/admin/permissions` | Frontend renderer only | Role and access controls |
| Compliance | `/admin/compliance` | Frontend renderer only | Compliance workspace |
| Users | `/admin/users` | Frontend renderer only | User directory |
| User detail | `/admin/users/[id]` | Frontend renderer only | Dynamic user detail |
| Operations | `/admin/operations` | Frontend renderer only | Ops dashboard |

## Implementation notes

1. **Frontend-only surface:** Every route in this file is handled by client-friendly page renderers. There is no dedicated backend API contract surfaced in this catch-all dispatcher.
2. **Data source:** Several routes rely on local mock data and shared UI components from the frontend app.
3. **If you want a stricter API contract:** the next step would be to split these into dedicated route modules or add server-side data hooks for the dynamic sections.

## Source of truth

- Current route implementation: [frontend/app/[...slug]/page.tsx](../frontend/app/[...slug]/page.tsx)
