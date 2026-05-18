# Orders Domain

Order management and processing system.

## Features
- Cart management
- Checkout flow
- Order tracking
- Inventory reservation
- Refunds
- Order status management
- Internal financial ledger
- Transactional consistency

## Order States
- PENDING
- PROCESSING
- PAID
- SHIPPED
- DELIVERED
- FAILED
- REFUNDED
- CANCELLED

## Entities
- Order
- OrderItem
- Payment
- Refund
- LedgerEntry

## Services
- OrderService - Order management
- CheckoutService - Checkout flow
- RefundService - Refund processing
- LedgerService - Financial ledger management

## API
- POST /checkout
- GET /orders
- GET /orders/:id
- POST /orders/:id/refund
- POST /payments/webhook
