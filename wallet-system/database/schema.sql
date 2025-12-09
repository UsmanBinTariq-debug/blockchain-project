-- Create tables for the crypto wallet system

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    cnic VARCHAR(20) UNIQUE,
    wallet_id VARCHAR(64) UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    encrypted_private_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_verified BOOLEAN DEFAULT FALSE
);

-- Wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    wallet_address VARCHAR(64) UNIQUE NOT NULL,
    balance_cache DECIMAL(20,8) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT NOW(),
    zakat_deducted_this_month BOOLEAN DEFAULT FALSE
);

-- UTXO table (Unspent Transaction Outputs)
CREATE TABLE IF NOT EXISTS utxos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_hash VARCHAR(64) NOT NULL,
    output_index INTEGER NOT NULL,
    wallet_address VARCHAR(64) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    is_spent BOOLEAN DEFAULT FALSE,
    spent_in_transaction VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(transaction_hash, output_index)
);

-- Blocks table
CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_index INTEGER UNIQUE NOT NULL,
    timestamp BIGINT NOT NULL,
    previous_hash VARCHAR(64) NOT NULL,
    hash VARCHAR(64) UNIQUE NOT NULL,
    nonce BIGINT NOT NULL,
    merkle_root VARCHAR(64),
    difficulty INTEGER DEFAULT 4,
    mined_by VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_hash VARCHAR(64) UNIQUE NOT NULL,
    block_hash VARCHAR(64) REFERENCES blocks(hash),
    sender_wallet VARCHAR(64) NOT NULL,
    receiver_wallet VARCHAR(64) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) DEFAULT 0,
    note TEXT,
    signature TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    transaction_type VARCHAR(20) DEFAULT 'transfer'
);

-- Zakat Transactions table
CREATE TABLE IF NOT EXISTS zakat_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(64) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    zakat_percentage DECIMAL(5,2) DEFAULT 2.5,
    transaction_hash VARCHAR(64) REFERENCES transactions(transaction_hash),
    month_year VARCHAR(7) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- System Logs table
CREATE TABLE IF NOT EXISTS system_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    log_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    wallet_address VARCHAR(64),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Beneficiaries table
CREATE TABLE IF NOT EXISTS beneficiaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    beneficiary_wallet_id VARCHAR(64) NOT NULL,
    nickname VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, beneficiary_wallet_id)
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_wallet_id ON users(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_address ON wallets(wallet_address);
CREATE INDEX IF NOT EXISTS idx_transactions_sender ON transactions(sender_wallet);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver ON transactions(receiver_wallet);
CREATE INDEX IF NOT EXISTS idx_transactions_hash ON transactions(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_blocks_index ON blocks(block_index);
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks(hash);
CREATE INDEX IF NOT EXISTS idx_utxos_wallet ON utxos(wallet_address);
CREATE INDEX IF NOT EXISTS idx_utxos_spent ON utxos(is_spent);
CREATE INDEX IF NOT EXISTS idx_zakat_wallet ON zakat_transactions(wallet_address);
CREATE INDEX IF NOT EXISTS idx_beneficiaries_user ON beneficiaries(user_id);
