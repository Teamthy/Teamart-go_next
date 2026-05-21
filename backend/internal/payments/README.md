# Payment Engine Architecture

## Overview

The Teamart Payment Engine is a comprehensive, multi-gateway payment processing system designed for marketplace platforms. It handles payments, refunds, escrow, split payments, merchant wallets, and payouts with support for global payment providers.

## Core Systems

### 1. Payment Processing (`service.go`)
- **Payment Intents**: Create and manage payment intents across multiple gateways
- **Payment Processing**: Process payments with full error handling
- **Payment Methods**: Save and manage multiple payment methods per user
- **Payment Status**: Track real-time payment status

### 2. Escrow System (`escrow.go`) - CRITICAL FOR MARKETPLACE
Protects buyers, sellers, and platform reputation:

```
Flow:
1. Buyer places order → Money enters escrow (held)
2. Order is delivered → Funds released to seller's wallet
3. If dispute → Funds held pending resolution
4. If cancellation → Funds refunded to buyer
```

**Features:**
- Automatic release on delivery confirmation
- Dispute management with resolution outcomes
- Support for split distribution (buyer/seller split)
- Audit trail for all escrow transactions

### 3. Split Payments Engine (`split_payments.go`)
Distributes payments to multiple recipients:

```
Example: $100 order
- Seller: $90
- Platform commission: $7
- Affiliate: $2
- Livestream host: $1
```

**Supports:**
- Platform fee/commission
- Affiliate commissions
- Creator/influencer commissions
- Livestream host revenue share
- Dynamic percentage-based or fixed amount splits

### 4. Merchant Wallet System (`merchant_wallet.go`)
Seller account management:

```
Balance Components:
- Total Balance: All available funds
- Pending Balance: In escrow, awaiting release
- Available Balance: Ready for withdrawal
```

**Features:**
- Credit sales and refunds
- Commission tracking
- Wallet transactions audit trail
- Automatic settlement workflows
- Support for multi-currency

### 5. Payout Engine (`payout_engine.go`)
Seller payment distribution:

**Supported Methods:**
- Bank Transfer
- Mobile Money (Airtel, MTN, Vodafone, etc.)
- PayPal
- Cryptocurrency
- Wallet Transfer

**Features:**
- Scheduled payouts (daily, weekly, bi-weekly, monthly)
- Instant payouts
- Minimum amount thresholds
- Automatic retry on failure
- Payout history and tracking

### 6. Refund Engine (`refund_engine.go`)
Handle all refund scenarios:

**Refund Reasons:**
- Customer request
- Product defective
- Product not received
- Not as described
- Duplicate charge
- Fraud
- Cancellation
- Return

**Workflow:**
1. Initiate refund (pending)
2. Approve refund (optional review)
3. Process refund (submit to gateway)
4. Confirm completion

### 7. Webhook Handler (`webhook_handler.go`)
Process payment gateway webhooks:

**Supported Providers:**
- Stripe
- Paystack
- Flutterwave
- Custom providers

**Webhook Events Handled:**
- Payment succeeded/failed
- Charge refunded
- Dispute created/resolved
- Payout completed/failed

### 8. Reconciliation (`reconciliation.go`)
Financial reconciliation and audit:

**Features:**
- Transaction matching
- Discrepancy identification and resolution
- Automatic reconciliation reports
- Pattern detection for fraud
- Full audit trail

## Payment Gateway Implementations

### Stripe (`stripe.go`)
- **API Version**: Stripe v1
- **Key Features**: 3D Secure, Webhooks, Instant payouts
- **Regions**: Global (primary for US/EU)
- **Currencies**: 135+ currencies

### Paystack (`paystack.go`)
- **Focus**: African payments
- **Supported Currencies**: NGN, GHS, KES, ZAR, UGX
- **Key Features**: Mobile money, Bank transfers, Webhooks
- **Regions**: Nigeria (primary), Ghana, Kenya, etc.

### Flutterwave (`flutterwave.go`)
- **Focus**: African and global payments
- **Supported Methods**: Cards, Wallets, Bank transfers
- **Key Features**: Payment links, Webhooks, Settlement APIs
- **Regions**: 40+ African countries, Global

### Alternative Payment Methods

#### Apple Pay (`alternative_gateways.go`)
- Server-side token validation and decryption
- Integration with primary payment processor
- Supports all order amounts

#### Google Pay (`alternative_gateways.go`)
- Server-side token handling
- Encrypted payload processing
- Mobile and web support

#### Mobile Money (`alternative_gateways.go`)
- Airtel Money
- MTN Mobile Money
- Vodafone Cash
- Orange Money
- Equitel

#### Cryptocurrency (`alternative_gateways.go`)
- Bitcoin
- Ethereum
- TRON
- Binance Chain
- Ripple

#### Bank Transfer (`alternative_gateways.go`)
- Virtual account numbers
- Direct bank transfers
- Automated verification

## Data Models

### Core Payment Models
- **PaymentMethod**: Saved payment instruments
- **PaymentIntent**: Payment request to gateway
- **PaymentTransaction**: Confirmed payment record
- **Refund**: Refund request and tracking
- **Wallet**: User/Seller account balance

### Marketplace Models
- **EscrowAccount**: Funds held for order protection
- **SplitPayment**: Multi-recipient payment distribution
- **MerchantWallet**: Seller earnings and balance
- **Payout**: Seller withdrawal request
- **PayoutSchedule**: Automatic payment scheduling

### Operations Models
- **WebhookLog**: Gateway webhook audit trail
- **PaymentReconciliation**: Financial reconciliation record
- **ReconciliationDiscrepancy**: Transaction mismatches

## Usage Examples

### Creating a Payment

```go
// Create payment intent
input := &CreatePaymentIntentInput{
    OrderID: 12345,
    UserID: 1,
    Amount: 99.99,
    Currency: "USD",
    Provider: "stripe", // or "paystack", "flutterwave"
}

result, err := paymentService.CreatePaymentIntent(ctx, input)
if err != nil {
    // Handle error
}

// Return payment intent to client
```

### Escrow Flow

```go
// 1. Create escrow when order is placed
escrow, err := escrowManager.CreateEscrow(ctx, orderID, buyerID, sellerID, totalAmount)

// 2. Release escrow when order is delivered
escrow, err := escrowManager.ReleaseEscrow(ctx, escrowID)
// → Funds automatically added to seller's wallet

// 3. Or refund if order is cancelled
escrow, err := escrowManager.RefundEscrow(ctx, escrowID)
// → Funds refunded to buyer
```

### Split Payment Distribution

```go
// Create split for $100 order
input := &CreateSplitPaymentInput{
    OrderID: 12345,
    TotalAmount: 100.0,
    Currency: "USD",
    Splits: []*SplitLine{
        {
            RecipientID: 1,        // Seller
            RecipientType: "seller",
            Amount: 90.0,
        },
        {
            RecipientID: 0,        // Platform
            RecipientType: "platform",
            Amount: 7.0,
        },
        {
            RecipientID: 2,        // Affiliate
            RecipientType: "affiliate",
            Amount: 2.0,
        },
        {
            RecipientID: 3,        // Livestream host
            RecipientType: "livestream_host",
            Amount: 1.0,
        },
    },
}

split, err := splitProcessor.CreateSplitPayment(ctx, input)
```

### Merchant Wallet Operations

```go
// Check balance
total, pending, available, err := walletManager.GetMerchantBalance(ctx, sellerID)

// Request payout
payout, err := walletManager.RequestPayout(ctx, sellerID, walletID, 1000.0)

// Approve payout
payout, err := walletManager.ApprovePayout(ctx, payoutID, adminUserID)

// Process payout
payout, err := walletManager.ProcessPayout(ctx, payoutID)
```

### Refund Processing

```go
// Initiate refund
refund, err := refundProcessor.InitiateRefund(ctx, orderID, "customer_request", requestedByUserID)

// Approve refund
refund, err := refundProcessor.ApproveRefund(ctx, refundID, adminUserID)

// Process refund
result, err := refundProcessor.ProcessRefund(ctx, refundID)
```

## Database Schema

All tables are created in migration: `005_payment_infrastructure.sql`

**Key Tables:**
- `payment_methods` - Saved payment instruments
- `payment_intents` - Payment requests
- `payment_transactions` - Payment records
- `refunds` - Refund tracking
- `wallet_transactions` - User wallet audit
- `escrow_accounts` - Escrow holds
- `escrow_disputes` - Dispute resolution
- `user_wallets` - User account balances
- `payouts` - Seller withdrawals
- `payment_gateway_config` - Gateway credentials (encrypted)
- `webhook_logs` - Webhook audit trail
- `transaction_taxes` - Tax records
- `financial_reconciliation` - Reconciliation reports

## Environment Variables Required

```
# Stripe
STRIPE_PUBLIC_KEY=pk_test_...
STRIPE_SECRET_KEY=sk_test_...

# Paystack
PAYSTACK_PUBLIC_KEY=pk_test_...
PAYSTACK_SECRET_KEY=sk_test_...

# Flutterwave
FLUTTERWAVE_PUBLIC_KEY=pk_test_...
FLUTTERWAVE_SECRET_KEY=sk_test_...

# Apple Pay
APPLE_PAY_MERCHANT_ID=merchant.com.yourapp

# Google Pay
GOOGLE_PAY_MERCHANT_ID=merchant_id_here

# Mobile Money
MOBILE_MONEY_API_KEY=your_api_key

# Webhook Secrets
STRIPE_WEBHOOK_SECRET=whsec_...
PAYSTACK_WEBHOOK_SECRET=...
FLUTTERWAVE_WEBHOOK_SECRET=...
```

## Security Considerations

1. **PCI Compliance**: Never store full credit card numbers
2. **Token Storage**: Always use provider tokens (Stripe tokens, etc.)
3. **Webhook Verification**: Always verify webhook signatures
4. **Encryption**: Gateway credentials stored encrypted in DB
5. **Rate Limiting**: Implement per-user payment rate limits
6. **Fraud Detection**: Integrate risk analysis (available in auth module)
7. **3D Secure**: Support for 3D Secure/Secure Authentication

## Testing

```bash
# Run tests
go test ./backend/internal/payments/...

# Test Stripe gateway
go test ./backend/internal/payments/ -run TestStripe

# Test Paystack gateway
go test ./backend/internal/payments/ -run TestPaystack

# Test reconciliation
go test ./backend/internal/payments/ -run TestReconciliation
```

## Monitoring & Alerts

Monitor these metrics:
- Payment success rate by gateway
- Average payment processing time
- Escrow hold time and release lag
- Payout completion time
- Refund processing time
- Reconciliation discrepancies
- Failed webhooks

## Future Enhancements

1. **Advanced Features**
   - Payment plans/subscriptions
   - Buy now, pay later (BNPL)
   - Installment payments
   - Dynamic pricing/currency conversion
   - Multi-tenant support

2. **Risk Management**
   - Chargeback handling
   - Fraud prevention rules
   - Velocity checks
   - Geographic restrictions

3. **Analytics**
   - Payment analytics dashboard
   - Revenue reporting
   - Settlement reporting
   - Tax compliance reporting

4. **Additional Gateways**
   - Square
   - Adyen
   - 2Checkout
   - Wise (international transfers)

## Support

For issues or questions about the payment engine:
1. Check this documentation
2. Review unit tests for usage examples
3. Check gateway-specific API documentation
4. Review migration SQL for database schema

## References

- [Stripe Documentation](https://stripe.com/docs)
- [Paystack Documentation](https://paystack.com/docs)
- [Flutterwave Documentation](https://developer.flutterwave.com/docs)
- PCI DSS Compliance Guide
- Best Practices for Payment Processing
