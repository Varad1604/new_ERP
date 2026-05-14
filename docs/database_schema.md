# Enterprise ERP Database Schema & ER Strategy

This document outlines the core database schema strategy, focusing on modularity, referential integrity, and institutional-grade financial correctness.

## 1. Core Modules

### 1.1 Identity & Access Management (IAM)

*   `users`: id, email, password_hash, mfa_secret, status, last_login.
*   `roles`: id, name, description, hierarchy_level.
*   `permissions`: id, resource, action (e.g., resource: 'ledger', action: 'write').
*   `role_permissions`: role_id, permission_id.
*   `user_roles`: user_id, role_id, department_id.
*   `departments`: id, name, parent_department_id.

### 1.2 Enterprise Operations (HR & Payroll)

*   `employees`: id, user_id, first_name, last_name, hire_date, salary_tier_id.
*   `attendance`: id, employee_id, date, clock_in, clock_out, status.
*   `payroll_runs`: id, period_start, period_end, status, total_amount, approved_by.
*   `payslips`: id, payroll_run_id, employee_id, base_pay, deductions, net_pay.

## 2. Institutional-Grade Financial Schema

The financial core strictly adheres to Double-Entry Accounting principles. Transactions are atomic and immutable.

### 2.1 The Chart of Accounts

*   `accounts`
    *   `id`: UUID (Primary Key)
    *   `code`: String (e.g., "1000", "2000")
    *   `name`: String (e.g., "Cash", "Accounts Receivable")
    *   `type`: Enum ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')
    *   `currency`: String (e.g., "USD")
    *   `is_active`: Boolean

### 2.2 Transactions & Ledgers

To ensure atomic double-entry logic, every financial event creates a `journal_entry` containing at least two `ledger_lines`.

*   `journal_entries` (The Header)
    *   `id`: UUID (Primary Key)
    *   `reference_number`: String (Unique, e.g., "INV-2023-001")
    *   `date`: Timestamp
    *   `description`: String
    *   `status`: Enum ('DRAFT', 'POSTED', 'VOIDED')
    *   `created_by`: UUID (FK to users)

*   `ledger_lines` (The Detail)
    *   `id`: UUID (Primary Key)
    *   `journal_entry_id`: UUID (FK to journal_entries)
    *   `account_id`: UUID (FK to accounts)
    *   `amount`: Decimal (Positive for Debit, Negative for Credit. Sum of all lines in an entry MUST equal 0).
    *   `cost_center_id`: UUID (Optional, for departmental tracking)

## 3. Data Integrity & Scalability

*   **Immutability:** Once a `journal_entry` is 'POSTED', it cannot be updated or deleted. Corrections require a new, reversing `journal_entry`.
*   **Audit Trails:** Every table modification triggers a database trigger/event that logs the change, user, and timestamp to a centralized `audit_logs` table or Elasticsearch index.
*   **Decimals:** All financial figures are stored as `DECIMAL(19,4)` or equivalent precise types. NEVER float/double to avoid rounding errors.
