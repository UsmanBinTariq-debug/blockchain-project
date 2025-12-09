# Architecture Decision Records

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
│                     Vercel (CDN + Serverless)                    │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  Pages: Login, Register, Dashboard, Wallet, Send, Explorer │  │
│  │  Components: Shared, Dashboard, Transactions, BlockExplorer│  │
│  │  State: Zustand (Auth, Wallet, UI)                         │  │
│  │  Styling: Tailwind CSS + Dark Mode                         │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────┬──────────────────────────────────┘
                               │ HTTPS
                               ↓
┌─────────────────────────────────────────────────────────────────┐
│                    Backend (Go) - Fly.io                         │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Gin HTTP Server                           │ │
│  │  ┌─────────────────────────────────────────────────────────┐│ │
│  │  │ API Routes & Handlers                                   ││ │
│  │  │ ├─ Auth (Register, Login, Verify OTP)                  ││ │
│  │  │ ├─ Wallet (Create, Balance, UTXOs)                     ││ │
│  │  │ ├─ Transactions (Send, Validate, History)              ││ │
│  │  │ └─ Blockchain (Mining, Validation, Explorer)           ││ │
│  │  └─────────────────────────────────────────────────────────┘│ │
│  │  ┌─────────────────────────────────────────────────────────┐│ │
│  │  │ Services Layer                                          ││ │
│  │  │ ├─ WalletService (Balance, Create, Validate)           ││ │
│  │  │ ├─ TransactionService (History, Status)                ││ │
│  │  │ ├─ ZakatService (Deduction, Reports)                   ││ │
│  │  │ └─ MiningService (PoW, Block Creation)                 ││ │
│  │  └─────────────────────────────────────────────────────────┘│ │
│  │  ┌─────────────────────────────────────────────────────────┐│ │
│  │  │ Blockchain Module                                       ││ │
│  │  │ ├─ Block (Structure, Hashing)                          ││ │
│  │  │ ├─ Transaction (Validation, UTXO)                      ││ │
│  │  │ ├─ ProofOfWork (Mining, Difficulty)                    ││ │
│  │  │ └─ Blockchain (Chain, UTXO Management)                 ││ │
│  │  └─────────────────────────────────────────────────────────┘│ │
│  │  ┌─────────────────────────────────────────────────────────┐│ │
│  │  │ Crypto Module                                           ││ │
│  │  │ ├─ Keys (RSA-2048, Encryption)                         ││ │
│  │  │ ├─ Signatures (ECDSA, Verification)                    ││ │
│  │  │ └─ Hashing (SHA-256, OTP)                              ││ │
│  │  └─────────────────────────────────────────────────────────┘│ │
│  │  ┌─────────────────────────────────────────────────────────┐│ │
│  │  │ Database Layer                                          ││ │
│  │  │ ├─ Models (User, Wallet, Transaction, Block, UTXO)    ││ │
│  │  │ ├─ PostgreSQL Driver (lib/pq)                          ││ │
│  │  │ └─ Connection Pooling (25 max, 5 idle)                 ││ │
│  │  └─────────────────────────────────────────────────────────┘│ │
│  └─────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────┘
                               │ TCP/PostgreSQL
                               ↓
┌─────────────────────────────────────────────────────────────────┐
│                  Supabase PostgreSQL Database                    │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │ Tables: Users, Wallets, Transactions, Blocks, UTXOs        │ │
│  │         ZakatTransactions, SystemLogs, Beneficiaries       │ │
│  │ Indexes: Email, WalletID, TransactionHash, BlockIndex      │ │
│  │ Foreign Keys: Relational integrity                         │ │
│  │ Backups: Automatic daily backups                           │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Key Architectural Decisions

### 1. UTXO Model over Account Model

**Decision**: Use UTXO (Unspent Transaction Output) model instead of account balance model.

**Rationale**:
- More transparent and auditable
- Matches Bitcoin/blockchain best practices
- Prevents double-spending naturally
- Enables parallel transaction processing
- Better for fee market dynamics

**Trade-offs**:
- More complex balance calculation
- Requires careful UTXO management
- More database queries

**Implementation**:
- Balance calculated as: `Sum of unspent UTXOs for wallet`
- Each transaction creates outputs (UTXOs)
- Spending marks UTXOs as used
- Caching for performance optimization

---

### 2. Proof-of-Work with Adjustable Difficulty

**Decision**: Implement PoW mining with dynamic difficulty adjustment.

**Rationale**:
- Secure and decentralized consensus
- Prevents spam and DOS attacks
- Adjusting difficulty maintains target mining time
- Proven in production systems

**Trade-offs**:
- Computationally expensive
- Variable block time without adjustment
- Resource intensive

**Implementation**:
- Target: 4 leading zeros (difficulty starts at 4)
- Adjusts ±1 based on mining time vs 1-minute target
- Formula: Hash < 2^(256-difficulty)
- Stored in database for historical tracking

---

### 3. Serverless PostgreSQL (Supabase)

**Decision**: Use Supabase's managed PostgreSQL instead of self-hosted database.

**Rationale**:
- Serverless scalability
- Automated backups and point-in-time recovery
- Built-in authentication and row-level security
- Real-time capabilities (future)
- Cost-effective for variable load

**Trade-offs**:
- Vendor lock-in
- Limited control over configuration
- Network latency considerations

**Implementation**:
- Connection pooling (PgBouncer)
- 25 max connections, 5 idle
- Automatic SSL/TLS
- Indexed columns for query performance

---

### 4. Go for Backend (Gin Framework)

**Decision**: Use Go with Gin web framework for the backend.

**Rationale**:
- High performance and concurrency
- Simple and clean syntax
- Strong standard library
- Excellent for distributed systems
- Easy deployment (single binary)

**Trade-offs**:
- Different paradigm from JavaScript
- Smaller ecosystem than Python/Node.js
- Verbose error handling

**Implementation**:
- Gin for HTTP routing and middleware
- Go's context for cancellation and timeouts
- Goroutines for concurrent processing
- Standard library for crypto operations

---

### 5. React + TypeScript Frontend

**Decision**: Use React with TypeScript and Tailwind CSS for the UI.

**Rationale**:
- React ecosystem maturity
- TypeScript for type safety
- Tailwind for rapid UI development
- Large community and resources
- Component reusability

**Trade-offs**:
- JavaScript bundle size
- Complexity for simple apps
- Learning curve

**Implementation**:
- React Router for navigation
- Zustand for state management
- React Hook Form for forms
- Vite for fast development and builds

---

### 6. JWT with OTP for Authentication

**Decision**: Use JWT tokens combined with OTP email verification.

**Rationale**:
- Stateless authentication
- Secure email verification
- User confirmation of login
- Compatible with distributed systems
- Can be revoked/refreshed

**Trade-offs**:
- OTP delivery latency
- Email dependency
- Token management complexity

**Implementation**:
- JWT signed with HS256
- OTP sent via Supabase Auth
- 15-minute OTP expiry
- Token refresh mechanism

---

### 7. RSA-2048 for Encryption and ECDSA for Signatures

**Decision**: Use RSA-2048 for key pairs and ECDSA for digital signatures.

**Rationale**:
- RSA-2048: Industry standard encryption
- ECDSA: More efficient for signatures
- Both production-proven
- Good balance of security and performance

**Trade-offs**:
- Larger key sizes than modern alternatives
- Computation cost for operations
- Compatibility with all systems

**Implementation**:
- Private keys encrypted with AES-256-GCM
- SHA-256 hashing for signatures
- Base64 encoding for storage
- Signature verification before transaction acceptance

---

### 8. GitHub Actions for CI/CD

**Decision**: Implement GitHub Actions for automated testing and deployment.

**Rationale**:
- Native GitHub integration
- Free for public repositories
- No additional infrastructure
- Familiar workflow syntax
- Good documentation

**Trade-offs**:
- Limited to GitHub
- Slower than self-hosted runners
- Limited concurrency

**Implementation**:
- Test workflows on PR
- Build and lint checks
- Deploy to Fly.io and Vercel on main branch
- Automated secret management

---

### 9. Vercel for Frontend, Fly.io for Backend

**Decision**: Separate deployments: Vercel for React app, Fly.io for Go backend.

**Rationale**:
- Vercel: Optimized for React, automatic builds from Git
- Fly.io: Excellent Go support, lightweight VMs
- Separation of concerns
- Independent scaling
- Easy to manage and monitor

**Trade-offs**:
- Cross-origin concerns (solved with CORS)
- Two separate deployments
- Potential latency between services

**Implementation**:
- Vercel: Deploy from `frontend/` directory
- Fly.io: Deploy from root with `fly.toml`
- Environment variables via dashboard
- Automatic SSL certificates

---

## Database Design Patterns

### 1. User Accounts
- UUID primary keys for distributed systems
- Unique email constraint for authentication
- CNIC for KYC compliance
- Encrypted private key storage
- OTP for email verification

### 2. Wallet Management
- One wallet per user (can extend to multiple)
- Cached balance for quick retrieval
- Zakat tracking for monthly deduction
- Last updated timestamp for staleness

### 3. Transaction Model
- Immutable transaction records
- Signature storage for audit trails
- Status tracking (pending, confirmed, failed)
- Transaction type classification
- Block hash reference for confirmed transactions

### 4. UTXO Tracking
- Unique constraint on (tx_hash, output_index)
- Is_spent boolean for quick filtering
- Spent_in_transaction for tracing
- Wallet address index for balance calculations

### 5. Block Storage
- Sequential block indices for integrity
- Hash uniqueness constraint
- Previous hash for chain validation
- Merkle root for transaction verification
- Nonce for PoW validation

---

## Error Handling Strategy

### Frontend
1. **Input Validation**: Form validation before submission
2. **API Errors**: Extract and display user-friendly messages
3. **Network Errors**: Retry logic with exponential backoff
4. **Auth Errors**: Redirect to login on 401
5. **Error Boundaries**: Catch React errors gracefully

### Backend
1. **Request Validation**: Validate all inputs immediately
2. **Business Logic Errors**: Return specific error codes
3. **Database Errors**: Log internally, return generic message
4. **Recovery**: Graceful degradation where possible
5. **Logging**: Structured logging for debugging

---

## Security Considerations

### Frontend
- ✅ XSS protection via React escaping
- ✅ CSRF tokens in forms
- ✅ Secure token storage (localStorage with caution)
- ✅ HTTPS only in production
- ✅ Content Security Policy headers

### Backend
- ✅ Input validation and sanitization
- ✅ SQL injection prevention (parameterized queries)
- ✅ Rate limiting on sensitive endpoints
- ✅ CORS validation
- ✅ JWT signature verification
- ✅ Private key encryption at rest
- ✅ Transaction signature verification
- ✅ HTTPS/TLS enforcement

### Database
- ✅ Encrypted connections (SSL/TLS)
- ✅ Row-level security (future enhancement)
- ✅ Audit logging
- ✅ Regular backups
- ✅ Access control via Supabase

---

## Performance Optimization

### Database
- Connection pooling reduces overhead
- Indexes on frequently queried columns
- Denormalized balance cache
- Efficient pagination
- Query optimization with EXPLAIN

### Backend
- Goroutine-based concurrency
- Caching layer for balance calculations
- Batch processing for zakat
- Efficient JSON serialization
- Gzip compression

### Frontend
- Code splitting with React Router
- Lazy loading of components
- Image optimization
- CSS-in-JS optimization
- Minimal re-renders with Zustand

---

## Monitoring and Observability

### Logging
- Structured JSON logging
- Log levels: debug, info, warn, error
- Request/response logging
- Performance metrics

### Metrics
- Request latency
- Error rates
- Database query performance
- Blockchain metrics (mining time, difficulty)

### Alerting
- Critical error alerts
- Deployment notifications
- Performance degradation alerts
- System health checks

---

## Future Enhancements

1. **Sharding**: Horizontal scaling for high TPS
2. **Layer 2 Solutions**: Off-chain transactions
3. **Multi-signature Wallets**: Require multiple approvals
4. **Hardware Wallet Integration**: Support ledger/trezor
5. **Staking**: Proof-of-Stake alternative
6. **Smart Contracts**: Programmable transactions
7. **DEX Integration**: Direct asset trading
8. **Mobile Apps**: React Native applications
9. **WebSocket Real-time**: Live balance updates
10. **Graph QL**: API query flexibility

---

**Last Updated**: January 2024
**Version**: 1.0.0
