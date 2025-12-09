package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database manages all database operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*Database, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// CreateUser creates a new user
func (d *Database) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, full_name, cnic, wallet_id, public_key, encrypted_private_key, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return d.db.QueryRowContext(ctx, query,
		user.Email, user.FullName, user.CNIC, user.WalletID,
		user.PublicKey, user.EncryptedPrivateKey, user.IsVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetUserByEmail retrieves a user by email
func (d *Database) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, full_name, cnic, wallet_id, public_key, encrypted_private_key, is_verified, created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &User{}
	err := d.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CNIC, &user.WalletID,
		&user.PublicKey, &user.EncryptedPrivateKey, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByWalletID retrieves a user by wallet ID
func (d *Database) GetUserByWalletID(ctx context.Context, walletID string) (*User, error) {
	query := `
		SELECT id, email, full_name, cnic, wallet_id, public_key, encrypted_private_key, is_verified, created_at, updated_at
		FROM users WHERE wallet_id = $1
	`

	user := &User{}
	err := d.db.QueryRowContext(ctx, query, walletID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CNIC, &user.WalletID,
		&user.PublicKey, &user.EncryptedPrivateKey, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, email, full_name, cnic, wallet_id, public_key, encrypted_private_key, is_verified, created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &User{}
	err := d.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.CNIC, &user.WalletID,
		&user.PublicKey, &user.EncryptedPrivateKey, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// UpdateUserVerification updates user verification status
func (d *Database) UpdateUserVerification(ctx context.Context, userID string, isVerified bool) error {
	query := `UPDATE users SET is_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := d.db.ExecContext(ctx, query, isVerified, userID)
	return err
}

// CreateWallet creates a new wallet
func (d *Database) CreateWallet(ctx context.Context, wallet *Wallet) error {
	query := `
		INSERT INTO wallets (user_id, wallet_address, balance_cache, last_updated)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	return d.db.QueryRowContext(ctx, query,
		wallet.UserID, wallet.WalletAddress, wallet.BalanceCache,
	).Scan(&wallet.ID)
}

// GetWalletsByUserID retrieves all wallets for a user
func (d *Database) GetWalletsByUserID(ctx context.Context, userID string) ([]*Wallet, error) {
	query := `
		SELECT id, user_id, wallet_address, balance_cache, last_updated, zakat_deducted_this_month
		FROM wallets WHERE user_id = $1
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*Wallet
	for rows.Next() {
		wallet := &Wallet{}
		err := rows.Scan(
			&wallet.ID, &wallet.UserID, &wallet.WalletAddress, &wallet.BalanceCache, &wallet.LastUpdated, &wallet.ZakatDeducted,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, rows.Err()
}

// GetWalletByAddress retrieves a wallet by address
func (d *Database) GetWalletByAddress(ctx context.Context, address string) (*Wallet, error) {
	query := `
		SELECT id, user_id, wallet_address, balance_cache, last_updated, zakat_deducted_this_month
		FROM wallets WHERE wallet_address = $1
	`

	wallet := &Wallet{}
	err := d.db.QueryRowContext(ctx, query, address).Scan(
		&wallet.ID, &wallet.UserID, &wallet.WalletAddress, &wallet.BalanceCache, &wallet.LastUpdated, &wallet.ZakatDeducted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return wallet, err
}

// UpdateWalletBalance updates the cached balance of a wallet
func (d *Database) UpdateWalletBalance(ctx context.Context, walletAddress string, balance float64) error {
	query := `UPDATE wallets SET balance_cache = $1, last_updated = NOW() WHERE wallet_address = $2`
	_, err := d.db.ExecContext(ctx, query, balance, walletAddress)
	return err
}

// CreateTransaction creates a new transaction record
func (d *Database) CreateTransaction(ctx context.Context, tx *Transaction) error {
	query := `
		INSERT INTO transactions (transaction_hash, sender_wallet, receiver_wallet, amount, fee, note, signature, status, transaction_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		tx.TransactionHash, tx.SenderWallet, tx.ReceiverWallet, tx.Amount,
		tx.Fee, tx.Note, tx.Signature, tx.Status, tx.TransactionType,
	).Scan(&tx.ID, &tx.CreatedAt)
}

// GetTransactionByHash retrieves a transaction by hash
func (d *Database) GetTransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	query := `
		SELECT id, transaction_hash, block_hash, sender_wallet, receiver_wallet, amount, fee, note, signature, status, created_at, transaction_type
		FROM transactions WHERE transaction_hash = $1
	`

	tx := &Transaction{}
	err := d.db.QueryRowContext(ctx, query, hash).Scan(
		&tx.ID, &tx.TransactionHash, &tx.BlockHash, &tx.SenderWallet, &tx.ReceiverWallet,
		&tx.Amount, &tx.Fee, &tx.Note, &tx.Signature, &tx.Status, &tx.CreatedAt, &tx.TransactionType,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tx, err
}

// GetTransactionsByWallet retrieves all transactions for a wallet
func (d *Database) GetTransactionsByWallet(ctx context.Context, walletAddress string, limit int, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, transaction_hash, block_hash, sender_wallet, receiver_wallet, amount, fee, note, signature, status, created_at, transaction_type
		FROM transactions
		WHERE sender_wallet = $1 OR receiver_wallet = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.db.QueryContext(ctx, query, walletAddress, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		tx := &Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.TransactionHash, &tx.BlockHash, &tx.SenderWallet, &tx.ReceiverWallet,
			&tx.Amount, &tx.Fee, &tx.Note, &tx.Signature, &tx.Status, &tx.CreatedAt, &tx.TransactionType,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// CreateBlock creates a new block record
func (d *Database) CreateBlock(ctx context.Context, block *Block) error {
	query := `
		INSERT INTO blocks (block_index, timestamp, previous_hash, hash, nonce, merkle_root, difficulty, mined_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		block.BlockIndex, block.Timestamp, block.PreviousHash, block.Hash,
		block.Nonce, block.MerkleRoot, block.Difficulty, block.MinedBy,
	).Scan(&block.ID, &block.CreatedAt)
}

// GetBlockByHash retrieves a block by hash
func (d *Database) GetBlockByHash(ctx context.Context, hash string) (*Block, error) {
	query := `
		SELECT id, block_index, timestamp, previous_hash, hash, nonce, merkle_root, difficulty, mined_by, created_at
		FROM blocks WHERE hash = $1
	`

	block := &Block{}
	err := d.db.QueryRowContext(ctx, query, hash).Scan(
		&block.ID, &block.BlockIndex, &block.Timestamp, &block.PreviousHash, &block.Hash,
		&block.Nonce, &block.MerkleRoot, &block.Difficulty, &block.MinedBy, &block.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return block, err
}

// CreateUTXO creates a new UTXO
func (d *Database) CreateUTXO(ctx context.Context, utxo *UTXO) error {
	query := `
		INSERT INTO utxos (transaction_hash, output_index, wallet_address, amount, is_spent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		utxo.TransactionHash, utxo.OutputIndex, utxo.WalletAddress, utxo.Amount, utxo.IsSpent,
	).Scan(&utxo.ID, &utxo.CreatedAt)
}

// GetUTXOsByWallet retrieves all UTXOs for a wallet
func (d *Database) GetUTXOsByWallet(ctx context.Context, walletAddress string) ([]*UTXO, error) {
	query := `
		SELECT id, transaction_hash, output_index, wallet_address, amount, is_spent, spent_in_transaction, created_at
		FROM utxos WHERE wallet_address = $1 AND is_spent = false
	`

	rows, err := d.db.QueryContext(ctx, query, walletAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var utxos []*UTXO
	for rows.Next() {
		utxo := &UTXO{}
		var spentIn sql.NullString
		err := rows.Scan(
			&utxo.ID, &utxo.TransactionHash, &utxo.OutputIndex, &utxo.WalletAddress,
			&utxo.Amount, &utxo.IsSpent, &spentIn, &utxo.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if spentIn.Valid {
			s := spentIn.String
			utxo.SpentInTransaction = &s
		} else {
			utxo.SpentInTransaction = nil
		}
		utxos = append(utxos, utxo)
	}

	return utxos, rows.Err()
}

// MarkUTXOAsSpent marks a UTXO as spent
func (d *Database) MarkUTXOAsSpent(ctx context.Context, txHash string, outputIndex int, spentInTx string) error {
	query := `
		UPDATE utxos SET is_spent = true, spent_in_transaction = $1
		WHERE transaction_hash = $2 AND output_index = $3
	`
	_, err := d.db.ExecContext(ctx, query, spentInTx, txHash, outputIndex)
	return err
}

// CreateZakatTransaction creates a zakat transaction record
func (d *Database) CreateZakatTransaction(ctx context.Context, zt *ZakatTransaction) error {
	query := `
		INSERT INTO zakat_transactions (wallet_address, amount, zakat_percentage, transaction_hash, month_year)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		zt.WalletAddress, zt.Amount, zt.ZakatPercentage, zt.TransactionHash, zt.MonthYear,
	).Scan(&zt.ID, &zt.CreatedAt)
}

// CreateSystemLog creates a system log entry
func (d *Database) CreateSystemLog(ctx context.Context, log *SystemLog) error {
	query := `
		INSERT INTO system_logs (log_type, message, wallet_address, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		log.LogType, log.Message, log.WalletAddress, log.IPAddress, log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
}

// GetSystemLogs retrieves system logs with optional filtering
func (d *Database) GetSystemLogs(ctx context.Context, logType string, limit int, offset int) ([]SystemLog, error) {
	query := `
		SELECT id, log_type, message, wallet_address, ip_address, user_agent, created_at
		FROM system_logs
	`
	args := []interface{}{}
	argIndex := 1

	if logType != "" && logType != "ALL" {
		query += fmt.Sprintf(" WHERE log_type = $%d", argIndex)
		args = append(args, logType)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []SystemLog
	for rows.Next() {
		var log SystemLog
		if err := rows.Scan(&log.ID, &log.LogType, &log.Message, &log.WalletAddress, &log.IPAddress, &log.UserAgent, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetSystemLogStats retrieves statistics about system logs
func (d *Database) GetSystemLogStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_logs,
			COUNT(CASE WHEN log_type = 'TRANSACTION_CREATED' THEN 1 END) as transaction_logs,
			COUNT(CASE WHEN log_type = 'BLOCK_MINED' THEN 1 END) as block_logs,
			COUNT(CASE WHEN log_type = 'ZAKAT_DEDUCTED' THEN 1 END) as zakat_logs,
			COUNT(CASE WHEN log_type = 'ERROR' THEN 1 END) as error_logs,
			COUNT(CASE WHEN log_type = 'AUTH' THEN 1 END) as auth_logs
		FROM system_logs
	`

	var stats map[string]interface{} = make(map[string]interface{})
	var totalLogs, transactionLogs, blockLogs, zakatLogs, errorLogs, authLogs int

	err := d.db.QueryRowContext(ctx, query).Scan(
		&totalLogs, &transactionLogs, &blockLogs, &zakatLogs, &errorLogs, &authLogs,
	)
	if err != nil {
		return nil, err
	}

	stats["total_logs"] = totalLogs
	stats["transaction_logs"] = transactionLogs
	stats["block_logs"] = blockLogs
	stats["zakat_logs"] = zakatLogs
	stats["error_logs"] = errorLogs
	stats["auth_logs"] = authLogs

	return stats, nil
}

// CreateBeneficiary creates a new beneficiary
func (d *Database) CreateBeneficiary(ctx context.Context, beneficiary *Beneficiary) error {
	query := `
		INSERT INTO beneficiaries (user_id, beneficiary_wallet_id, nickname)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return d.db.QueryRowContext(ctx, query,
		beneficiary.UserID, beneficiary.BeneficiaryWalletID, beneficiary.Nickname,
	).Scan(&beneficiary.ID, &beneficiary.CreatedAt)
}

// GetBeneficiariesByUserID retrieves all beneficiaries for a user
func (d *Database) GetBeneficiariesByUserID(ctx context.Context, userID string) ([]*Beneficiary, error) {
	query := `
		SELECT id, user_id, beneficiary_wallet_id, nickname, created_at
		FROM beneficiaries WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beneficiaries []*Beneficiary
	for rows.Next() {
		b := &Beneficiary{}
		err := rows.Scan(&b.ID, &b.UserID, &b.BeneficiaryWalletID, &b.Nickname, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		beneficiaries = append(beneficiaries, b)
	}

	return beneficiaries, rows.Err()
}

// GetBlocks retrieves blocks with pagination
func (d *Database) GetBlocks(ctx context.Context, limit int, offset int) ([]*Block, error) {
	query := `
		SELECT id, block_index, timestamp, previous_hash, hash, nonce, merkle_root, difficulty, mined_by, created_at
		FROM blocks
		ORDER BY block_index DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*Block
	for rows.Next() {
		block := &Block{}
		err := rows.Scan(
			&block.ID, &block.BlockIndex, &block.Timestamp, &block.PreviousHash, &block.Hash,
			&block.Nonce, &block.MerkleRoot, &block.Difficulty, &block.MinedBy, &block.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, rows.Err()
}

// GetTransactionsByBlockHash retrieves transactions for a block
func (d *Database) GetTransactionsByBlockHash(ctx context.Context, blockHash string, limit int, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, transaction_hash, block_hash, sender_wallet, receiver_wallet, amount, fee, signature, status, transaction_type, note, created_at
		FROM transactions
		WHERE block_hash = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.db.QueryContext(ctx, query, blockHash, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		tx := &Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.TransactionHash, &tx.BlockHash, &tx.SenderWallet, &tx.ReceiverWallet,
			&tx.Amount, &tx.Fee, &tx.Signature, &tx.Status, &tx.TransactionType, &tx.Note, &tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// GetTransactionsByStatus retrieves transactions by status
func (d *Database) GetTransactionsByStatus(ctx context.Context, status string, limit int) ([]*Transaction, error) {
	query := `
		SELECT id, transaction_hash, block_hash, sender_wallet, receiver_wallet, amount, fee, signature, status, transaction_type, note, created_at
		FROM transactions
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := d.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		tx := &Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.TransactionHash, &tx.BlockHash, &tx.SenderWallet, &tx.ReceiverWallet,
			&tx.Amount, &tx.Fee, &tx.Signature, &tx.Status, &tx.TransactionType, &tx.Note, &tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// UpdateTransactionStatus updates a transaction's status and block hash
func (d *Database) UpdateTransactionStatus(ctx context.Context, txHash string, status string, blockHash string) error {
	query := `
		UPDATE transactions
		SET status = $1, block_hash = $2, updated_at = NOW()
		WHERE transaction_hash = $3
	`

	_, err := d.db.ExecContext(ctx, query, status, blockHash, txHash)
	return err
}

// Ping checks the database connection
func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
