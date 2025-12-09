# Deployment Guide

## Prerequisites

- Git account (GitHub)
- Supabase account (https://supabase.com)
- Fly.io account (https://fly.io)
- Vercel account (https://vercel.com)
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

## Database Setup (Supabase)

### Step 1: Create Supabase Project

1. Go to https://supabase.com
2. Click "New Project"
3. Enter project details:
   - Project name: `crypto-wallet-system`
   - Password: Generate strong password
   - Region: Select closest to your users
4. Wait for project initialization

### Step 2: Run Database Schema

1. Go to SQL Editor in Supabase dashboard
2. Click "New Query"
3. Copy and paste contents of `database/schema.sql`
4. Click "Run"
5. Verify all tables are created

### Step 3: Get Connection Details

1. Go to Project Settings → Database
2. Copy:
   - Connection string (URI)
   - Host
   - Database name
   - User
   - Password

Example connection string:
```
postgresql://postgres:[PASSWORD]@db.[PROJECT_ID].supabase.co:5432/postgres
```

### Step 4: Set Up Backups

1. Go to Project Settings → Backups
2. Set backup frequency to Daily
3. Set retention to 30 days

## Backend Deployment (Fly.io)

### Step 1: Install Flyctl

```bash
curl -L https://fly.io/install.sh | sh
```

Add to PATH:
```bash
export PATH="$PATH:$HOME/.fly/bin"
```

### Step 2: Create Fly Account

```bash
flyctl auth signup
# or
flyctl auth login
```

### Step 3: Prepare Backend

```bash
cd backend

# Create .env file with production values
cp .env.example .env

# Update .env with:
# - DATABASE_URL (from Supabase)
# - JWT_SECRET (generate with: openssl rand -hex 32)
# - ZAKAT_POOL_WALLET (create a special wallet)
# - Other environment variables
```

### Step 4: Launch Application

```bash
flyctl launch
```

When prompted:
```
? App Name: crypto-wallet-backend
? Organization: (select your organization)
? Region: (select closest region)
? Would you like to set up a Postgresql database: No
? Would you like to set up an upstash redis: No
```

### Step 5: Set Secrets

```bash
flyctl secrets set \
  DATABASE_URL="postgresql://..." \
  JWT_SECRET="generated-secret" \
  ZAKAT_POOL_WALLET="wallet-address" \
  SUPABASE_URL="https://project.supabase.co" \
  SUPABASE_ANON_KEY="anon-key" \
  SUPABASE_SERVICE_ROLE_KEY="service-role-key" \
  CORS_ALLOWED_ORIGINS="https://your-frontend.vercel.app" \
  NODE_ENV="production" \
  LOG_LEVEL="info"
```

### Step 6: Deploy

```bash
flyctl deploy

# Monitor deployment
flyctl logs
```

### Step 7: Get Backend URL

```bash
flyctl info

# Your app URL will be displayed
# Example: https://crypto-wallet-backend.fly.dev
```

## Frontend Deployment (Vercel)

### Step 1: Push to GitHub

```bash
git add .
git commit -m "Initial commit"
git push origin main
```

### Step 2: Import to Vercel

1. Go to https://vercel.com
2. Click "New Project"
3. Select your GitHub repository
4. Choose "Frontend" as root directory: `frontend`

### Step 3: Set Environment Variables

In Vercel dashboard:

1. Go to Settings → Environment Variables
2. Add:
   ```
   VITE_API_URL=https://crypto-wallet-backend.fly.dev/api
   VITE_SUPABASE_URL=https://project.supabase.co
   VITE_SUPABASE_ANON_KEY=your-anon-key
   ```

### Step 4: Deploy

1. Click "Deploy"
2. Wait for build to complete
3. Get your frontend URL: `https://your-project.vercel.app`

### Step 5: Update Backend CORS

Back in Fly.io:

```bash
cd backend

flyctl secrets set \
  CORS_ALLOWED_ORIGINS="https://your-frontend.vercel.app,http://localhost:5173"

flyctl deploy
```

## GitHub Actions Setup

### Step 1: Generate API Tokens

**For Fly.io:**
```bash
flyctl auth token
# Copy the token
```

**For Vercel:**
1. Go to Account Settings → Tokens
2. Create new token with scope: `Full Account Access`

### Step 2: Add Secrets to GitHub

1. Go to repository → Settings → Secrets and variables → Actions
2. Add secrets:
   - `FLY_API_TOKEN`: Fly.io token
   - `VERCEL_TOKEN`: Vercel token
   - `VERCEL_ORG_ID`: From Vercel account
   - `VERCEL_PROJECT_ID`: From Vercel project

To get Vercel IDs:
```bash
vercel projects ls
vercel env
```

### Step 3: GitHub Actions Workflow

The `.github/workflows/deploy.yml` is already configured. On push to main:
1. Tests run
2. Backend builds
3. Frontend builds
4. Deploys to Fly.io and Vercel if tests pass

## Monitoring and Maintenance

### Backend Monitoring (Fly.io)

```bash
# View logs
flyctl logs

# Check app status
flyctl status

# View metrics
flyctl metrics

# SSH into machine (if needed)
flyctl ssh console
```

### Database Monitoring (Supabase)

1. Go to Supabase dashboard
2. Check Statistics for usage
3. View Logs for query performance
4. Monitor Storage for size

### Frontend Monitoring (Vercel)

1. Go to Vercel dashboard
2. View Analytics for traffic
3. Check Deployments for history
4. Monitor Functions (serverless endpoints)

## Backup and Recovery

### Database Backups

**Automatic Backups:**
- Supabase: Daily backups retained for 30 days
- Access from: Project Settings → Backups

**Manual Backup:**
```bash
# Dump database
pg_dump postgresql://user:password@host:port/dbname > backup.sql

# Restore database
psql postgresql://user:password@host:port/dbname < backup.sql
```

### Code Backups

```bash
# Create git backup
git tag backup-v1.0.0
git push origin backup-v1.0.0

# Clone from backup
git clone --branch backup-v1.0.0 repository.git
```

## Scaling Considerations

### Database Scaling

1. **Read Replicas**: Add replicas for read-heavy workloads
2. **Caching**: Implement Redis for frequently accessed data
3. **Sharding**: Partition data by wallet address (future)

### Backend Scaling

1. **Horizontal Scaling**: Add more Fly.io machines
   ```bash
   flyctl scale count -a crypto-wallet-backend 2
   ```

2. **Vertical Scaling**: Increase machine size
   ```bash
   flyctl machine update <machine-id> --memory 512
   ```

3. **Load Balancing**: Fly.io handles automatically

### Frontend Scaling

Vercel automatically scales based on traffic. No action needed.

## Cost Optimization

### Reduce Costs

1. **Database**: 
   - Use appropriate indexes
   - Archive old data
   - Reduce connection count

2. **Backend**:
   - Use smallest machine size that works
   - Scale down during off-hours
   - Optimize queries

3. **Frontend**:
   - Enable compression
   - Optimize images
   - Use CDN (Vercel provides)

## Security Hardening

### TLS/SSL

Both Fly.io and Vercel provide automatic HTTPS. Verify:

```bash
# Check certificate
curl -I https://crypto-wallet-backend.fly.dev
# Should show: certificate verify ok

curl -I https://your-frontend.vercel.app
# Should show: certificate verify ok
```

### Environment Variables

1. Never commit `.env` files
2. Use Fly.io `flyctl secrets`
3. Use Vercel environment variables UI
4. Rotate secrets regularly

### Database Security

1. Enable SSL connections:
   ```
   sslmode=require
   ```

2. Use strong passwords
3. Limit connection sources
4. Enable audit logging

### API Security

1. Enable rate limiting
2. Add request signing
3. Use API keys for service-to-service
4. Monitor for suspicious activity

## Troubleshooting

### Backend Won't Deploy

```bash
# Check build logs
flyctl logs

# Verify dependencies
go mod tidy
go mod verify

# Local test
go run ./cmd/server/main.go

# Check fly.toml
flyctl config show
```

### Database Connection Issues

```bash
# Test connection string
psql "postgresql://user:pass@host:port/db"

# Check firewall rules
flyctl ips list

# Verify Supabase settings
# Go to: Project Settings → Database → Connection Info
```

### Frontend Not Loading

```bash
# Check build logs in Vercel
vercel logs

# Verify environment variables
vercel env ls

# Test locally
npm run build
npm run preview

# Check API connectivity
curl https://crypto-wallet-backend.fly.dev/health
```

### CORS Errors

```bash
# Update CORS in backend
flyctl secrets set CORS_ALLOWED_ORIGINS="https://new-domain.com"

flyctl deploy

# Test
curl -H "Origin: https://new-domain.com" \
     -H "Access-Control-Request-Method: POST" \
     https://crypto-wallet-backend.fly.dev/api/health
```

## Rollback Procedure

### Backend Rollback

```bash
# View deployment history
flyctl releases list

# Rollback to previous version
flyctl releases rollback <version>

# Or redeploy from git
git revert <commit-hash>
git push origin main
flyctl deploy
```

### Frontend Rollback

1. Go to Vercel dashboard
2. Click "Deployments"
3. Find previous working deployment
4. Click three dots → "Promote to Production"

## Disaster Recovery

### Complete System Recovery

1. **Database**: Restore from Supabase backup
   ```bash
   # From Supabase dashboard: Backups → Restore
   ```

2. **Backend**: Redeploy from GitHub
   ```bash
   flyctl deploy
   ```

3. **Frontend**: Automatic through Vercel

### Data Recovery

```bash
# Export data from backup
pg_dump backup-database.sql > recovery.sql

# Check data integrity
psql < recovery.sql
```

## Production Checklist

- [ ] Database backups enabled
- [ ] Monitoring and alerting configured
- [ ] API rate limiting enabled
- [ ] HTTPS/TLS certificates valid
- [ ] Environment variables set correctly
- [ ] Error logging configured
- [ ] Performance optimization done
- [ ] Security audit completed
- [ ] Documentation updated
- [ ] Team access configured
- [ ] Incident response plan ready
- [ ] Load testing completed

## Support Resources

- Fly.io Docs: https://fly.io/docs
- Vercel Docs: https://vercel.com/docs
- Supabase Docs: https://supabase.com/docs
- Go Docs: https://golang.org/doc
- React Docs: https://react.dev

---

**Last Updated**: January 2024
**Version**: 1.0.0
