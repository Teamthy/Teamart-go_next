# Authentication Domain

Authentication state machine implementation with multi-step onboarding flow.

## Features
- XState-based state machine
- OTP verification
- JWT token management
- Email verification
- Progressive onboarding

## Entities
- User authentication states
- OTP tokens
- JWT tokens
- Onboarding progress

## Services
- AuthService - Core authentication logic
- OTPService - One-time password management
- TokenService - JWT token generation and validation

## API
- POST /auth/signup
- POST /auth/verify-email
- POST /auth/resend-otp
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout
- GET /auth/me
- GET /auth/onboarding-state
