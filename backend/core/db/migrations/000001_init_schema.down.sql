DROP INDEX IF EXISTS idx_ledger_lines_tenant;
DROP INDEX IF EXISTS idx_journal_entries_tenant;
DROP INDEX IF EXISTS idx_accounts_tenant;
DROP INDEX IF EXISTS idx_users_tenant;
DROP INDEX IF EXISTS idx_roles_tenant;

DROP TABLE IF EXISTS ledger_lines;
DROP TABLE IF EXISTS journal_entries;
DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS journal_status;
DROP TYPE IF EXISTS account_type;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS tenants;
