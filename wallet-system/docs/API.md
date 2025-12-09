# API Documentation

## Base URL
```
Development: http://localhost:8080/api
Production: https://your-backend.fly.dev/api
```

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer {jwt_token}
```

## Response Format

All API responses follow this format:

### Success Response
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "error": "Error message",
  "message": "Detailed error message",
  "code": "ERROR_CODE"
}
```

---

## Auth Endpoints

### Register User
**POST** `/auth/register`

Creates a new user account with wallet.

Request:
```json
{
  "email": "user@example.com",
  "full_name": "John Doe",
  "cnic": "12345-1234567-1",
  "password": "SecurePass123!"
}
```

Response:
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "wallet_id": "64-char-hex",
    "wallet_address": "64-char-hex"
  }
}
```

Error Cases:
- `INVALID_EMAIL` - Invalid email format
- `WEAK_PASSWORD` - Password doesn't meet requirements
- `INVALID_CNIC` - Invalid CNIC format
- `EMAIL_EXISTS` - Email already registered
- `USER_CREATE_ERROR` - Failed to create user
- `WALLET_CREATE_ERROR` - Failed to create wallet

---

### Login
**POST** `/auth/login`

Authenticates user and sends OTP to email.

Request:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

Response:
```json
{
  "status": "success",
  "message": "OTP sent to email. Please verify.",
  "data": {
    "user_id": "uuid",
    "message": "Check your email for OTP"
  }
}
```

Error Cases:
- `INVALID_CREDENTIALS` - Wrong email or password
- `OTP_ERROR` - Failed to send OTP
- `DB_ERROR` - Database error

---

### Verify OTP
**POST** `/auth/verify-otp`

Verifies OTP and returns JWT token.

Request:
```json
{
  "user_id": "uuid",
  "otp": "123456"
}
```

Response:
```json
{
  "status": "success",
  "message": "Email verified successfully",
  "data": {
    "token": "jwt-token"
  }
}
```

Error Cases:
- `INVALID_OTP` - Invalid OTP format
- `OTP_EXPIRED` - OTP has expired
- `DB_ERROR` - Database error

---

## Wallet Endpoints

### Get Wallet Profile
**GET** `/wallet/profile`

Requires authentication.

Response:
```json
{
  "status": "success",
  "message": "Wallet retrieved",
  "data": {
    "user_id": "uuid",
    "wallet_id": "64-char-hex",
    "wallet_address": "64-char-hex",
    "balance": 100.5,
    "last_updated": "2024-01-01T12:00:00Z"
  }
}
```

---

### Get Balance
**POST** `/wallet/balance`

Calculates balance from UTXOs for a wallet.

Request:
```json
{
  "wallet_address": "64-char-hex-address"
}
```

Response:
```json
{
  "status": "success",
  "message": "Balance retrieved",
  "data": {
    "wallet_address": "64-char-hex-address",
    "balance": 100.5,
    "utxo_count": 5
  }
}
```

Error Cases:
- `INVALID_WALLET` - Invalid wallet address format
- `BALANCE_ERROR` - Failed to calculate balance

---

### Get UTXOs
**GET** `/wallet/utxos?wallet_address=<address>`

Returns all unspent transaction outputs for a wallet.

Response:
```json
{
  "status": "success",
  "message": "UTXOs retrieved",
  "data": {
    "utxos": [
      {
        "id": "uuid",
        "transaction_hash": "64-char-hex",
        "output_index": 0,
        "wallet_address": "64-char-hex",
        "amount": 50.25,
        "is_spent": false,
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

## Transaction Endpoints

### Send Money
**POST** `/transaction/send`

Creates and validates a new transaction.

Request:
```json
{
  "sender_wallet": "64-char-hex",
  "receiver_wallet": "64-char-hex",
  "amount": 10.5,
  "fee": 0.001,
  "note": "Payment for services",
  "signature": "base64-encoded-signature"
}
```

Response:
```json
{
  "status": "success",
  "message": "Transaction created",
  "data": {
    "transaction_hash": "64-char-hex",
    "sender_wallet": "64-char-hex",
    "receiver_wallet": "64-char-hex",
    "amount": 10.5,
    "fee": 0.001,
    "status": "pending",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

Error Cases:
- `INVALID_WALLET` - Invalid sender or receiver wallet
- `INVALID_SIGNATURE` - Invalid digital signature
- `INSUFFICIENT_BALANCE` - Not enough balance
- `UTXO_ALREADY_SPENT` - UTXO has already been spent
- `INVALID_AMOUNT` - Invalid transaction amount

---

### Get Pending Transactions
**GET** `/transaction/pending`

Returns all pending transactions waiting to be mined.

Response:
```json
{
  "status": "success",
  "message": "Pending transactions retrieved",
  "data": {
    "transactions": [
      {
        "id": "uuid",
        "transaction_hash": "64-char-hex",
        "sender_wallet": "64-char-hex",
        "receiver_wallet": "64-char-hex",
        "amount": 10.5,
        "fee": 0.001,
        "status": "pending",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

### Get Transaction History
**GET** `/transaction/history?wallet_address=<address>&limit=10&offset=0`

Returns transaction history for a wallet.

Query Parameters:
- `wallet_address` (required): Wallet address
- `limit` (optional): Number of records (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

Response:
```json
{
  "status": "success",
  "message": "Transaction history retrieved",
  "data": {
    "total": 50,
    "transactions": [
      {
        "id": "uuid",
        "transaction_hash": "64-char-hex",
        "sender_wallet": "64-char-hex",
        "receiver_wallet": "64-char-hex",
        "amount": 10.5,
        "fee": 0.001,
        "status": "confirmed",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

## Blockchain Endpoints

### Get All Blocks
**GET** `/blockchain/blocks?limit=10&offset=0`

Returns all blocks in the blockchain.

Query Parameters:
- `limit` (optional): Number of blocks (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

Response:
```json
{
  "status": "success",
  "message": "Blocks retrieved",
  "data": {
    "total": 100,
    "blocks": [
      {
        "id": "uuid",
        "block_index": 0,
        "timestamp": 1704067200,
        "previous_hash": "0",
        "hash": "64-char-hex",
        "nonce": 1234,
        "merkle_root": "64-char-hex",
        "difficulty": 4,
        "mined_by": "64-char-hex",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

### Get Block by Hash
**GET** `/blockchain/blocks/<hash>`

Returns a specific block.

Response:
```json
{
  "status": "success",
  "message": "Block retrieved",
  "data": {
    "block": {
      "id": "uuid",
      "block_index": 0,
      "timestamp": 1704067200,
      "previous_hash": "0",
      "hash": "64-char-hex",
      "nonce": 1234,
      "merkle_root": "64-char-hex",
      "difficulty": 4,
      "mined_by": "64-char-hex",
      "transactions": []
    }
  }
}
```

---

### Get Latest Block
**GET** `/blockchain/latest`

Returns the latest block in the chain.

Response:
```json
{
  "status": "success",
  "message": "Latest block retrieved",
  "data": {
    "block": {
      "id": "uuid",
      "block_index": 100,
      "timestamp": 1704067200,
      "previous_hash": "64-char-hex",
      "hash": "64-char-hex",
      "nonce": 5678,
      "merkle_root": "64-char-hex",
      "difficulty": 5
    }
  }
}
```

---

### Mine Block
**POST** `/blockchain/mine`

Mines a new block with pending transactions.

Request:
```json
{
  "miner_address": "64-char-hex"
}
```

Response:
```json
{
  "status": "success",
  "message": "Block mined successfully",
  "data": {
    "block": {
      "id": "uuid",
      "block_index": 101,
      "timestamp": 1704067260,
      "previous_hash": "64-char-hex",
      "hash": "64-char-hex",
      "nonce": 9012,
      "merkle_root": "64-char-hex",
      "difficulty": 5,
      "mined_by": "64-char-hex"
    }
  }
}
```

---

## Beneficiary Endpoints

### Add Beneficiary
**POST** `/beneficiary/add`

Requires authentication.

Request:
```json
{
  "beneficiary_wallet_id": "64-char-hex",
  "nickname": "Friend John"
}
```

Response:
```json
{
  "status": "success",
  "message": "Beneficiary added",
  "data": {
    "id": "uuid",
    "beneficiary_wallet_id": "64-char-hex",
    "nickname": "Friend John"
  }
}
```

---

### Get Beneficiaries
**GET** `/beneficiary/list`

Requires authentication.

Response:
```json
{
  "status": "success",
  "message": "Beneficiaries retrieved",
  "data": {
    "beneficiaries": [
      {
        "id": "uuid",
        "beneficiary_wallet_id": "64-char-hex",
        "nickname": "Friend John",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

---

## Health Check

### Health Status
**GET** `/health`

Returns API health status.

Response:
```json
{
  "status": "ok"
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|------------|-------------|
| `INVALID_REQUEST` | 400 | Invalid request format |
| `INVALID_EMAIL` | 400 | Invalid email format |
| `WEAK_PASSWORD` | 400 | Password doesn't meet requirements |
| `INVALID_CNIC` | 400 | Invalid CNIC format |
| `INVALID_WALLET` | 400 | Invalid wallet address |
| `INVALID_AMOUNT` | 400 | Invalid transaction amount |
| `INVALID_OTP` | 400 | Invalid OTP format |
| `INVALID_SIGNATURE` | 401 | Invalid digital signature |
| `INSUFFICIENT_BALANCE` | 400 | Insufficient balance |
| `UTXO_ALREADY_SPENT` | 400 | UTXO has already been spent |
| `EMAIL_EXISTS` | 409 | Email already registered |
| `UNAUTHORIZED` | 401 | Unauthorized access |
| `NOT_FOUND` | 404 | Resource not found |
| `DB_ERROR` | 500 | Database error |
| `SERVER_ERROR` | 500 | Internal server error |

---

## Rate Limiting

- **Default Limit**: 100 requests per minute per IP
- **Auth Endpoints**: 5 requests per minute
- **Transaction Endpoints**: 10 requests per minute

---

## CORS

Allowed origins (configurable):
- `http://localhost:5173` (development)
- `https://your-frontend.vercel.app` (production)

---

## Pagination

Endpoints supporting pagination accept:
- `limit` (default: 10, max: 100)
- `offset` (default: 0)

---

## Sorting

Endpoints with sorting support:
- Sort by: `created_at`, `amount`, `status`
- Order: `asc` (default), `desc`

---

**Last Updated**: January 2024
**API Version**: 1.0.0
