# Crypto Wallet System - Quick Start Guide

## 🎯 Quick Start (5 minutes)

### Prerequisites
- Docker installed
- Git installed

### Option 1: Docker Compose (Easiest)

```bash
# Clone repository
git clone https://github.com/yourusername/crypto-wallet-system.git
cd wallet-system

# Start with docker-compose
docker-compose up
```

The app will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api
- Database: PostgreSQL on localhost:5432

### Option 2: Local Development

#### Backend Setup
```bash
cd backend

# Copy environment
cp .env.example .env

# Install dependencies
go mod download

# Run server
go run ./cmd/server/main.go
```

Backend running at: http://localhost:8080

#### Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

Frontend running at: http://localhost:5173

#### Database Setup
1. Create Supabase project: https://supabase.com
2. Run SQL from `database/schema.sql` in Supabase SQL editor
3. Update `backend/.env` with your DATABASE_URL

## 📚 Documentation

- [API Documentation](docs/API.md) - Complete API reference
- [Architecture Guide](docs/ARCHITECTURE.md) - System design and decisions
- [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment steps

## 🔐 Security Features

✅ RSA-2048 encryption for wallet keys
✅ SHA-256 hashing for blocks
✅ ECDSA digital signatures for transactions
✅ AES-256-GCM encryption for private keys
✅ JWT authentication with OTP verification
✅ Rate limiting on sensitive endpoints
✅ SQL injection prevention
✅ XSS protection

## 🚀 Features

### Blockchain
- ✅ Custom blockchain implementation
- ✅ UTXO model for balance tracking
- ✅ Proof-of-Work mining with adjustable difficulty
- ✅ Merkle tree for transaction verification
- ✅ Chain validation and integrity checks

### Wallet
- ✅ Automatic wallet creation
- ✅ Public/private key pair generation
- ✅ Encrypted key storage
- ✅ QR code for wallet address
- ✅ Multiple transaction history views
- ✅ Beneficiary management

### Transactions
- ✅ Digital signature verification
- ✅ UTXO-based transaction validation
- ✅ Double-spend prevention
- ✅ Transaction fee calculation
- ✅ Real-time balance updates
- ✅ Transaction status tracking

### Zakat
- ✅ Automatic 2.5% monthly deduction
- ✅ Zakat transaction history
- ✅ Monthly reports
- ✅ Scheduled processing on 1st of month

### UI/UX
- ✅ Modern React interface
- ✅ Dark/light mode support
- ✅ Responsive design
- ✅ Real-time updates
- ✅ Error handling and validation
- ✅ Transaction history and reports

## 📊 Technology Stack

**Backend**: Go 1.21+, Gin Framework
**Frontend**: React 18, TypeScript, Tailwind CSS
**Database**: PostgreSQL (via Supabase)
**Deployment**: Fly.io (backend), Vercel (frontend)

## 🧪 Testing

```bash
# Backend tests
cd backend
go test -v ./...

# Frontend tests
cd frontend
npm test
```

## 🔧 Troubleshooting

### Port Already in Use
```bash
# Backend (8080)
lsof -i :8080
kill -9 <PID>

# Frontend (5173)
lsof -i :5173
kill -9 <PID>
```

### Database Connection Error
```bash
# Test connection
psql "postgresql://user:pass@host:port/db"

# Update .env with correct DATABASE_URL
```

### Frontend Build Error
```bash
# Clear cache
rm -rf node_modules package-lock.json
npm install
npm run build
```

## 📞 Support

- GitHub Issues: [Create Issue](https://github.com/yourusername/crypto-wallet-system/issues)
- Email: support@cryptowallet.example.com
- Documentation: See docs/ directory

## 📜 License

MIT License - See LICENSE file

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🎉 Getting Started

1. **Register**: Create a new account with email verification
2. **View Wallet**: Check your wallet address and QR code
3. **Receive Funds**: Share your wallet address with others
4. **Send Money**: Send crypto to other users
5. **Track History**: View all transactions and block explorer
6. **Monitor Zakat**: Track automatic zakat deductions

## 🔐 Default Test Accounts

For development, use these test credentials:

| Email | Password |
|-------|----------|
| test@example.com | TestPass123! |
| user@example.com | UserPass123! |

## 📈 Production Checklist

- [ ] Update environment variables
- [ ] Set up database backups
- [ ] Configure monitoring and alerting
- [ ] Enable HTTPS
- [ ] Set up CI/CD pipeline
- [ ] Configure rate limiting
- [ ] Test disaster recovery
- [ ] Document runbooks
- [ ] Set up log aggregation
- [ ] Configure APM monitoring

---

**Last Updated**: January 2024
**Version**: 1.0.0
**Status**: Production Ready ✅
