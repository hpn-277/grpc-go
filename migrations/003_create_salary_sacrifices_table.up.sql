-- Create salary_sacrifices table
CREATE TABLE IF NOT EXISTS salary_sacrifices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    super_account_id UUID NOT NULL REFERENCES super_accounts(id) ON DELETE CASCADE,
    amount_cents BIGINT NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_amount_positive CHECK (amount_cents > 0),
    CONSTRAINT chk_frequency CHECK (frequency IN ('weekly', 'fortnightly', 'monthly'))
);

-- Create index on super_account_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_salary_sacrifices_super_account_id ON salary_sacrifices(super_account_id);

-- Create index on start_date for date range queries
CREATE INDEX IF NOT EXISTS idx_salary_sacrifices_start_date ON salary_sacrifices(start_date DESC);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_salary_sacrifices_created_at ON salary_sacrifices(created_at DESC);
