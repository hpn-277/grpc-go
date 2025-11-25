-- Drop indexes
DROP INDEX IF EXISTS idx_salary_sacrifices_created_at;
DROP INDEX IF EXISTS idx_salary_sacrifices_start_date;
DROP INDEX IF EXISTS idx_salary_sacrifices_super_account_id;

-- Drop salary_sacrifices table
DROP TABLE IF EXISTS salary_sacrifices CASCADE;
