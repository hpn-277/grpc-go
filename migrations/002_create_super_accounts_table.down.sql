-- Drop indexes
DROP INDEX IF EXISTS idx_super_accounts_user_account;
DROP INDEX IF EXISTS idx_super_accounts_created_at;
DROP INDEX IF EXISTS idx_super_accounts_user_id;

-- Drop super_accounts table
DROP TABLE IF EXISTS super_accounts CASCADE;
