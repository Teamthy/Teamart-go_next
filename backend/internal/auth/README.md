# Phase 2 — Authentication Engine & Identity Infrastructure

For a production-grade AI-native livestream commerce platform like TikTok Shop Seller Center, the authentication engine is not just “login/signup.”

It is a:

```text
distributed identity and access infrastructure
```

supporting:

* merchants
* creators
* customers
* admins
* support agents
* multi-device sessions
* livestream moderation
* payouts/KYC
* tenant isolation
* fraud prevention
* realtime authorization

---

# 🎯 PHASE 2 GOAL

Build:

```text
a production-grade identity platform
```

NOT:

```text
basic auth routes
```

---

# 🏗️ PHASE 2 ARCHITECTURE

```text
Client Apps
    ↓
API Gateway / Middleware
    ↓
Identity Service
    ↓
Session Engine
    ↓
RBAC + Permission Engine
    ↓
Security Engine
    ↓
Onboarding State Machine
    ↓
Redis + PostgreSQL
```

---

# 📦 PHASE 2 MODULES

Your auth engine should contain these domains:

```text
/internal/auth
    /identity
    /session
    /token
    /rbac
    /onboarding
    /security
    /otp
    /middleware
    /audit
    /device
    /kyc
```

---

# 1. IDENTITY DOMAIN

This is the core user identity infrastructure.

---

## Responsibilities

* account creation
* identity lifecycle
* email verification
* account status
* tenant ownership
* account recovery
* identity federation readiness

---

## Core Models

### User

```go
type User struct {
    ID              uuid.UUID
    Email           string
    PasswordHash    string

    Role            Role
    Status          AccountStatus

    EmailVerified   bool
    PhoneVerified   bool

    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

---

## Account Status

```text
PENDING
ACTIVE
SUSPENDED
BANNED
DELETED
```

---

## Required APIs

```http
POST /auth/signup
POST /auth/login
POST /auth/logout
GET  /auth/me
PATCH /auth/profile
POST /auth/forgot-password
POST /auth/reset-password
```

---

# 2. ONBOARDING STATE MACHINE

This is CRITICAL.

TikTok Shop-style onboarding is multi-step and resumable.

---

# Merchant Onboarding States

```text
START
EMAIL_ENTERED
EMAIL_VERIFIED
PROFILE_CREATED
BUSINESS_CREATED
STORE_CREATED
KYC_SUBMITTED
BANK_CONNECTED
LIVE_COMMERCE_ENABLED
COMPLETED
```

---

# Requirements

* persistent onboarding
* resumable flows
* step validation
* state transition rules
* audit logs
* progress tracking

---

## Database Table

```sql
merchant_onboarding_states
```

---

## State Transition Engine

```go
type OnboardingService interface {
    AdvanceState(ctx context.Context, userID uuid.UUID, nextState State) error
    GetCurrentState(ctx context.Context, userID uuid.UUID) (*State, error)
}
```

---

# 3. SESSION ENGINE

This is NOT optional.

You need centralized session management.

---

# Responsibilities

* multi-device sessions
* refresh token rotation
* session revocation
* device tracking
* suspicious login detection
* concurrent session limits

---

# Session Model

```go
type Session struct {
    ID              uuid.UUID
    UserID          uuid.UUID

    RefreshTokenID  string

    DeviceID        string
    DeviceName      string
    IPAddress       string
    UserAgent       string

    ExpiresAt       time.Time
    RevokedAt       *time.Time

    CreatedAt       time.Time
}
```

---

# Required Features

## Token Rotation

Every refresh:

```text
old refresh token → revoked
new refresh token → issued
```

---

## Session Revocation

Support:

```text
logout current device
logout all devices
admin-forced logout
```

---

## Redis Session Cache

```text
session:{id}
user_sessions:{user_id}
```

---

# 4. TOKEN ENGINE

Separate token logic from auth logic.

---

# Responsibilities

* JWT generation
* refresh token generation
* token validation
* token rotation
* signing key management
* token blacklisting

---

# JWT Claims

```go
type Claims struct {
    UserID     string
    Role       string
    SessionID  string
    TenantID   string

    jwt.RegisteredClaims
}
```

---

# Access Token Rules

```text
short-lived
5–15 mins
```

---

# Refresh Token Rules

```text
long-lived
stored securely
rotated
revocable
```

---

# 5. OTP ENGINE

For email verification + sensitive actions.

---

# Responsibilities

* OTP generation
* OTP delivery
* OTP expiration
* resend throttling
* verification limits

---

# Redis Keys

```text
otp:{email}
otp_attempts:{email}
otp_resend:{email}
```

---

# Security Rules

```text
6-digit OTP
10-minute expiry
max attempts
rate-limited resend
```

---

# APIs

```http
POST /auth/verify-email
POST /auth/resend-otp
POST /auth/verify-otp
```

---

# 6. RBAC ENGINE

Critical for multi-tenant systems.

---

# Roles

```text
SUPER_ADMIN
PLATFORM_ADMIN
MERCHANT
CREATOR
CUSTOMER
SUPPORT_AGENT
```

---

# Responsibilities

* role permissions
* tenant isolation
* ownership validation
* route protection
* stream moderation permissions

---

# Permission Model

```go
type Permission struct {
    Resource string
    Action   string
}
```

---

# Example

```text
products:create
products:update
streams:start
orders:refund
admin:moderate
```

---

# Required Middleware

```go
RequireAuth()
RequireRole()
RequirePermission()
RequireStoreOwnership()
```

---

# 7. SECURITY ENGINE

This is a major subsystem.

---

# Responsibilities

* password hashing
* rate limiting
* suspicious activity detection
* IP reputation
* CSRF protection
* brute-force protection
* token abuse prevention
* fraud signals

---

# Required Features

## Password Hashing

Use:

```text
bcrypt
```

---

## Login Rate Limits

```text
Redis throttling
```

---

## Abuse Prevention

Track:

* failed logins
* OTP abuse
* token replay
* suspicious refresh attempts

---

# Redis Keys

```text
login_attempts:{ip}
login_attempts:{email}
blocked_ip:{ip}
```

---

# 8. DEVICE MANAGEMENT ENGINE

TikTok Shop-style platforms require device visibility.

---

# Responsibilities

* trusted devices
* suspicious device detection
* login history
* active sessions

---

# Features

```text
list active devices
revoke device
new device alerts
device fingerprinting
```

---

# APIs

```http
GET    /sessions
DELETE /sessions/:id
DELETE /sessions/all
```

---

# 9. AUDIT LOGGING ENGINE

MANDATORY for production systems.

---

# Responsibilities

Track:

* auth events
* admin actions
* session revocations
* permission changes
* fraud events
* payout access

---

# Audit Event Example

```json
{
  "event": "user.login",
  "user_id": "uuid",
  "ip": "127.0.0.1",
  "device": "Chrome",
  "timestamp": "2026-05-18T20:00:00Z"
}
```

---

# Audit Tables

```text
audit_logs
security_events
login_attempts
```

---

# 10. KYC + COMPLIANCE ENGINE

Needed for merchants + creators.

---

# Responsibilities

* merchant verification
* creator verification
* document uploads
* payout eligibility
* sanctions/fraud checks

---

# States

```text
PENDING
UNDER_REVIEW
APPROVED
REJECTED
```

---

# APIs

```http
POST /kyc/submit
GET  /kyc/status
POST /kyc/documents
```

---

# 🧠 AUTH EVENT-DRIVEN ARCHITECTURE

Auth should publish Kafka events.

---

# Events

```text
auth.signup.completed
auth.email.verified
auth.login.success
auth.login.failed
auth.session.revoked
auth.password.changed
merchant.onboarding.completed
kyc.submitted
```

---

# Consumers

## Notifications

```text
send welcome email
send suspicious login alert
```

---

## Analytics

```text
track onboarding conversion
track login metrics
```

---

## Fraud Engine

```text
detect account abuse
detect fake merchants
```

---

# 🗄️ DATABASE TABLES

You should have:

```text
users
sessions
refresh_tokens
otp_codes
merchant_onboarding_states
roles
permissions
role_permissions
user_roles
audit_logs
security_events
devices
kyc_submissions
```

---

# 🔒 SECURITY REQUIREMENTS

MANDATORY:

* bcrypt hashing
* JWT rotation
* secure cookies
* CSRF protection
* Redis throttling
* rate limiting
* audit logging
* session revocation
* refresh token rotation
* tenant isolation
* webhook signature validation

---

# ⚡ REDIS RESPONSIBILITIES

Redis should handle:

```text
OTP cache
session cache
rate limiting
temporary onboarding state
blacklisted tokens
security throttling
```

---

# 🧩 PHASE 2 OUTPUT EXPECTATION

At the end of Phase 2, your platform should support:

✅ secure signup/login
✅ merchant onboarding state machine
✅ multi-device sessions
✅ refresh token rotation
✅ RBAC authorization
✅ session revocation
✅ Redis-backed OTP infrastructure
✅ audit logging
✅ fraud prevention foundation
✅ tenant-aware identity infrastructure
✅ production-grade middleware
✅ creator/admin/customer/merchant roles
✅ distributed auth event publishing
✅ secure identity lifecycle management

---

# 📁 RECOMMENDED FOLDER STRUCTURE

```text
/internal/auth
    /identity
        entity.go
        service.go
        repository.go
        handler.go

    /session
        entity.go
        service.go
        repository.go

    /token
        jwt.go
        refresh.go

    /otp
        service.go
        redis.go

    /rbac
        permissions.go
        middleware.go

    /security
        throttling.go
        abuse_detection.go

    /onboarding
        machine.go
        transitions.go

    /audit
        logger.go

    /device
        fingerprint.go

    /kyc
        service.go
```
