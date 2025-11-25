-- Create super_accounts table
CREATE TABLE IF NOT EXISTS super_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fund_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(100) NOT NULL,
    balance_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_super_accounts_user_id ON super_accounts(user_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_super_accounts_created_at ON super_accounts(created_at DESC);

-- Create unique constraint on account_number per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_super_accounts_user_account ON super_accounts(user_id, account_number);
