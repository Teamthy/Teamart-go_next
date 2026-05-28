export interface CartItem {
    id: string;
    name: string;
    price: string;
    qty: number;
    shipping?: string;
    refundStatus?: string;
}

export interface Order {
    id: string;
    customer: string;
    total: string;
    status: string;
    items: number;
    date: string;
    tracking?: string;
    paymentMethod?: string;
    refundStatus?: string;
    shipping?: string;
}

export interface Payout {
    id: string;
    amount: string;
    period: string;
    status: string;
}

export interface OrderStatusItem {
    label: string;
    value: string;
}

export interface CheckoutDetail {
    label: string;
    value: string;
}

export const cartItems: CartItem[] = [
    { id: "1", name: "Sustainable Canvas Tote", price: "$32", qty: 1, shipping: "Priority express", refundStatus: "None" },
    { id: "2", name: "Artist Collaboration Hoodie", price: "$74", qty: 2, shipping: "Standard shipping", refundStatus: "None" },
];

export const sellerOrders: Order[] = [
    { id: "O-5021", customer: "Mia R.", total: "$248", status: "Processing", items: 6, date: "May 20", tracking: "TRK-8821", paymentMethod: "Visa ending in 2412", refundStatus: "No refund", shipping: "Priority express" },
    { id: "O-5018", customer: "Nate K.", total: "$124", status: "Shipped", items: 3, date: "May 19", tracking: "TRK-8810", paymentMethod: "PayPal", refundStatus: "No refund", shipping: "Standard" },
    { id: "O-5015", customer: "Asha T.", total: "$398", status: "Pending", items: 9, date: "May 18", tracking: "Pending", paymentMethod: "Apple Pay", refundStatus: "Pending review", shipping: "Priority" },
    { id: "O-5012", customer: "Lina W.", total: "$76", status: "Delivered", items: 2, date: "May 17", tracking: "TRK-8794", paymentMethod: "Visa ending in 2412", refundStatus: "Refunded", shipping: "Standard" },
    { id: "O-5009", customer: "Theo P.", total: "$199", status: "Delivered", items: 5, date: "May 16", tracking: "TRK-8782", paymentMethod: "Mastercard ending in 4444", refundStatus: "No refund", shipping: "Express" },
    { id: "O-5004", customer: "Jade B.", total: "$112", status: "Canceled", items: 1, date: "May 14", tracking: "Canceled", paymentMethod: "Visa ending in 2412", refundStatus: "Canceled", shipping: "Not shipped" },
];

export const sellerPayouts: Payout[] = [
    { id: "P-308", amount: "$1,220", period: "May 11 - May 17", status: "Completed" },
    { id: "P-309", amount: "$1,750", period: "May 18 - May 24", status: "Pending" },
    { id: "P-310", amount: "$430", period: "May 25 - May 31", status: "Scheduled" },
];

export const orderTrackingSummary: OrderStatusItem[] = [
    { label: "Processing", value: "2 orders" },
    { label: "Shipped", value: "1 order" },
    { label: "Delivered", value: "2 orders" },
    { label: "Refunds", value: "1 pending" },
];

export const orderRefundNotes = [
    "Refund requested for O-5004 due to size mismatch",
    "Customer support has already contacted the buyer",
    "A replacement order is scheduled for the next drop window",
];

export const paymentMethods = [
    { label: "Visa ending in 2412", value: "Primary" },
    { label: "PayPal", value: "Secondary" },
    { label: "Apple Pay", value: "Fast checkout" },
];

export const checkoutSteps = [
    "Shipping details confirmed",
    "Payment method selected",
    "Review and confirm",
];

export const checkoutRecoverySteps = [
    "Payment details need a quick revisit",
    "Shipping details are still saved",
    "You can retry without losing your cart",
];

export const orderConfirmationItems = [
    "Order placed successfully",
    "Shipping confirmation will arrive shortly",
    "You can track your order from account history",
];

export const checkoutShippingDetails: CheckoutDetail[] = [
    { label: "Customer", value: "Mia Rivera" },
    { label: "Address", value: "42 Pink Lane, Apt 3, San Francisco, CA 94107" },
    { label: "Mode", value: "Priority express" },
    { label: "Notes", value: "Leave package at the front desk if unavailable" },
];

export const checkoutShippingNotes = [
    "Deliver between 9 AM and 5 PM",
    "Leave package at the front desk if unavailable",
    "SMS updates are enabled",
];

export const checkoutPaymentSummary: CheckoutDetail[] = [
    { label: "Method", value: "Visa ending in 2412" },
    { label: "Subtotal", value: "$106" },
    { label: "Shipping", value: "$8" },
    { label: "Total", value: "$114" },
];

export const checkoutConfidencePoints = [
    "Fraud checks passed",
    "Customer support available if delivery changes",
    "Order confirmation will be sent instantly",
];

export const checkoutReviewChecklist = [
    "Shipping confirmed and ready",
    "Payment method confirmed",
    "Order total verified",
    "Support contacts available",
];

export const productEditorChecklist = [
    "Confirm product title and description",
    "Review pricing and discount settings",
    "Verify highlight media and availability",
];

export const launchReadinessChecklist = [
    "Ready for creator bundle placement",
    "Shipping and promo settings synced",
    "Preview available in storefront theme",
];

export const checkoutRecoveryActions = [
    "Retry payment",
    "Edit cart",
];

export const orderInsights = [
    "Conversion is strongest on live drop days",
    "Priority express orders are trending upward",
    "Refund risk is lowest on creator bundle purchases",
];
