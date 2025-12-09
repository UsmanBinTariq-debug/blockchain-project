# Crypto Wallet System

A fully functional decentralized cryptocurrency wallet system with custom blockchain implementation, Proof-of-Work mining, digital signatures, and automatic Zakat deduction.

## Features

- ✅ Custom blockchain using UTXO model
- ✅ Proof-of-Work mining with SHA-256 hashing
- ✅ ECDSA digital signatures for transaction security
- ✅ Automatic 2.5% monthly Zakat deduction
- ✅ User authentication with OTP verification
- ✅ Real-time balance tracking via UTXO model
- ✅ Complete transaction validation system
- ✅ Block explorer interface
- ✅ Transaction history and reporting
- ✅ QR code wallet sharing
- ✅ Beneficiary management
- ✅ Dark/light mode support

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: Supabase PostgreSQL (Serverless)
- **Cryptography**: RSA 2048-bit keys, SHA-256, AES encryption
- **Authentication**: JWT tokens with OTP verification

### Frontend
- **Framework**: React 18 + TypeScript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite
- **State Management**: Zustand
- **Forms**: React Hook Form
- **Charts**: Recharts
- **QR Code**: qrcode.react

### Deployment
- **Backend**: Fly.io
- **Frontend**: Vercel
- **Database**: Supabase PostgreSQL
- **CI/CD**: GitHub Actions

## Project Structure

```
wallet-system/
├── backend/
│   ├── cmd/server/
│   │   └── main.go
│   ├── internal/
│   │   ├── blockchain/
│   │   ├── crypto/
│   │   ├── database/
│   │   ├── api/
│   │   ├── services/
│   │   └── utils/
│   ├── pkg/config/
│   ├── tests/
│   ├── go.mod
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── types/
│   │   ├── utils/
│   │   └── App.tsx
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
├── database/
│   ├── schema.sql
│   └── migrations/
├── docs/
│   ├── API.md
│   ├── ARCHITECTURE.md
│   └── DEPLOYMENT.md
└── README.md
```

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL (via Supabase)
- Git
- Make (for backend)

## Installation

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/crypto-wallet-system.git
cd wallet-system
```

### 2. Backend Setup

```bash
cd backend

# Copy environment variables
cp .env.example .env

# Install dependencies
go mod download

# Run server
go run ./cmd/server/main.go
```

### 3. Frontend Setup

```bash
cd frontend

# Copy environment variables
cp .env.example .env.development

# Install dependencies
npm install

# Start development server
npm run dev
```

### 4. Database Setup

1. Create a Supabase project at https://supabase.com
2. Run the schema in `database/schema.sql`:
   ```sql
   -- Copy and run all SQL from database/schema.sql in Supabase SQL editor
   ```

## API Documentation

### Authentication Endpoints

**POST /api/auth/register**
```json
{
  "email": "user@example.com",
  "full_name": "John Doe",
  "cnic": "12345-1234567-1",
  "password": "SecurePass123!"
}
```

**POST /api/auth/login**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**POST /api/auth/verify-otp**
```json
{
  "user_id": "uuid",
  "otp": "123456"
}
```

### Wallet Endpoints

**GET /api/wallet/profile** (Requires Auth)
- Returns: User wallet information

**POST /api/wallet/balance**
```json
{
  "wallet_address": "64-char-hex-address"
}
```
- Returns: Current wallet balance from UTXOs

### Transaction Endpoints

**POST /api/transaction/send**
```json
{
  "sender_wallet": "...",
  "receiver_wallet": "...",
  "amount": 10.5,
  "fee": 0.001,
  "note": "Payment",
  "signature": "..."
}
```

**GET /api/transaction/pending**
- Returns: List of pending transactions

**GET /api/blockchain/blocks**
- Returns: All blocks in the blockchain

**POST /api/blockchain/mine**
```json
{
  "miner_address": "..."
}
```

## Database Schema

### Users
- `id` (UUID)
- `email` (VARCHAR, UNIQUE)
- `full_name` (VARCHAR)
- `cnic` (VARCHAR, UNIQUE)
- `wallet_id` (VARCHAR, UNIQUE)
- `public_key` (TEXT)
- `encrypted_private_key` (TEXT)
- `is_verified` (BOOLEAN)
- `otp_code` (VARCHAR)
- `otp_expires_at` (TIMESTAMP)

### Wallets
- `id` (UUID)
- `user_id` (UUID, FK)
- `wallet_address` (VARCHAR, UNIQUE)
- `balance_cache` (DECIMAL)
- `zakat_deducted_this_month` (BOOLEAN)

### UTXOs (Unspent Transaction Outputs)
- `id` (UUID)
- `transaction_hash` (VARCHAR)
- `output_index` (INTEGER)
- `wallet_address` (VARCHAR)
- `amount` (DECIMAL)
- `is_spent` (BOOLEAN)
- `spent_in_transaction` (VARCHAR)

### Blocks
- `id` (UUID)
- `block_index` (INTEGER, UNIQUE)
- `timestamp` (BIGINT)
- `previous_hash` (VARCHAR)
- `hash` (VARCHAR, UNIQUE)
- `nonce` (BIGINT)
- `merkle_root` (VARCHAR)
- `difficulty` (INTEGER)
- `mined_by` (VARCHAR)

### Transactions
- `id` (UUID)
- `transaction_hash` (VARCHAR, UNIQUE)
- `block_hash` (VARCHAR, FK)
- `sender_wallet` (VARCHAR)
- `receiver_wallet` (VARCHAR)
- `amount` (DECIMAL)
- `fee` (DECIMAL)
- `signature` (TEXT)
- `status` (VARCHAR)
- `transaction_type` (VARCHAR)

### Zakat Transactions
- `id` (UUID)
- `wallet_address` (VARCHAR)
- `amount` (DECIMAL)
- `zakat_percentage` (DECIMAL)
- `month_year` (VARCHAR)

## Blockchain Implementation

### Block Structure
```go
type Block struct {
    Index        int64
    Timestamp    int64
    Transactions []Transaction
    PreviousHash string
    Nonce        int64
    Hash         string
    MerkleRoot   string
    Difficulty   int
    MinedBy      string
}
```

### Transaction Structure
```go
type Transaction struct {
    ID             string
    SenderWallet   string
    ReceiverWallet string
    Amount         float64
    Fee            float64
    Note           string
    Timestamp      int64
    Signature      string
    PublicKey      string
    UTXOInputs     []UTXO
    UTXOOutputs    []UTXO
    Status         string
}
```

### UTXO Model
Balance is calculated as the sum of unspent transaction outputs for a wallet:
```
Balance = Sum of all unspent UTXOs for wallet address
```

### Proof-of-Work
- Difficulty starts at 4 leading zeros
- Target: `hash < 2^(256-difficulty)`
- Difficulty adjusts based on mining time (target: 1 minute per block)

### Digital Signatures
- Algorithm: RSA-2048
- Hash: SHA-256
- Signature verification required for all transactions

### Zakat System
- Automatically deducts 2.5% of wallet balance on the 1st of each month
- Creates a zakat transaction record
- Transfers funds to `ZAKAT_POOL_WALLET`

## Deployment

### Deploy to Fly.io (Backend)

1. **Install Flyctl**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Authenticate**
   ```bash
   flyctl auth login
   ```

3. **Create App**
   ```bash
   cd backend
   flyctl app create crypto-wallet-backend
   ```

4. **Set Secrets**
   ```bash
   flyctl secrets set DATABASE_URL="postgresql://..."
   flyctl secrets set JWT_SECRET="your-secret"
   # Set other environment variables
   ```

5. **Deploy**
   ```bash
   flyctl deploy
   ```

### Deploy to Vercel (Frontend)

1. **Push to GitHub**
   ```bash
   git push origin main
   ```

2. **Connect to Vercel**
   - Go to https://vercel.com
   - Import project from GitHub
   - Set environment variables:
     - `VITE_API_URL`: Your Fly.io backend URL
     - `VITE_SUPABASE_URL`: Your Supabase URL
     - `VITE_SUPABASE_ANON_KEY`: Your Supabase key

3. **Deploy**
   - Vercel automatically deploys on push to main

## Testing

### Backend Tests
```bash
cd backend
go test -v ./...
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests
```bash
cd frontend
npm test
npm run test:coverage
```

## Security Considerations

1. **Private Keys**: Always encrypted before storage using AES-256
2. **Signatures**: All transactions must be digitally signed
3. **HTTPS Only**: All APIs use HTTPS in production
4. **JWT Tokens**: Implement token refresh and expiry
5. **Rate Limiting**: Implement on API endpoints
6. **Input Validation**: All inputs are validated
7. **XSS Protection**: Using React's built-in protections
8. **SQL Injection**: Using parameterized queries

## Development Guidelines

### Backend
- Follow Go code standards (gofmt, golint)
- Use middleware for cross-cutting concerns
- Implement proper error handling
- Write unit tests for business logic
- Use context for cancellation and timeouts

### Frontend
- Use TypeScript for type safety
- Component-based architecture
- Custom hooks for reusable logic
- Proper error boundaries
- Responsive design with Tailwind CSS

## Performance Optimization

### Backend
- Connection pooling (25 max open, 5 idle)
- Database query optimization with indexes
- Caching strategy for balance calculations
- Asynchronous task scheduling for Zakat

### Frontend
- Code splitting with React Router
- Lazy loading of components
- Image optimization
- Minification and compression
- Service Worker for offline support (future)

## Monitoring and Logging

### Backend
- Structured JSON logging
- Log levels: debug, info, warn, error
- Request/response logging middleware
- Performance metrics tracking
- Error tracking and reporting

### Frontend
- Error boundaries for crash handling
- Performance monitoring
- User interaction tracking
- Error reporting to backend

## Troubleshooting

### Database Connection Issues
```bash
# Test connection string
psql your_connection_string

# Check Supabase status
# Go to https://status.supabase.com
```

### Frontend Build Issues
```bash
# Clear node_modules
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Backend Build Issues
```bash
# Clear Go cache
go clean -cache
go mod tidy
go build ./cmd/server
```

## Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Commit changes: `git commit -am 'Add feature'`
3. Push to branch: `git push origin feature/your-feature`
4. Create Pull Request

## License

MIT License - See LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: https://github.com/yourusername/crypto-wallet-system/issues
- Email: support@cryptowallet.example.com

## Roadmap

- [ ] Mobile app (React Native)
- [ ] Multi-signature wallets
- [ ] Hardware wallet integration
- [ ] Automated market maker (AMM)
- [ ] Decentralized exchange (DEX)
- [ ] Advanced analytics dashboard
- [ ] Institutional features
- [ ] DeFi integrations

## Changelog

### v1.0.0 (Initial Release)
- Basic wallet creation and management
- Transaction validation and signing
- Block explorer interface
- Zakat calculation system
- Complete user authentication
- Dashboard and reporting

---

**Built with ❤️ for the crypto community**
