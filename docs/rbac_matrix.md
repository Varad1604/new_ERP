# Role-Based Access Control (RBAC) Matrix

Enterprise security requires fine-grained, dynamic, and hierarchical access control. This ERP uses a Policy-Based Access Control (PBAC) model built on top of standard RBAC.

## 1. Core Principles

*   **Principle of Least Privilege:** Users only have access to resources necessary for their specific job function.
*   **Separation of Duties (SoD):** A single user cannot complete a critical business process alone (e.g., the person who creates a vendor cannot approve payments to that vendor).
*   **Contextual Access:** Permissions can be conditional (e.g., Manager can approve expenses *up to $10,000*; CEO required for *> $10,000*).

## 2. Standard Roles Overview

| Role | Scope | Key Capabilities |
| :--- | :--- | :--- |
| **Super Admin** | Global | Full system access, infrastructure configuration, global user management. |
| **C-Level / Board** | Global | Read-only access to all high-level dashboards, financial forecasts, and company-wide analytics. |
| **Finance Controller**| Financial Domain | Full access to general ledger, journal entries, payroll approval, and audit trails. |
| **Finance Clerk** | Financial Domain | Can draft journal entries (requires approval to post), manage basic AP/AR. |
| **HR Manager** | HR Domain | Full access to employee records, payroll drafting, onboarding/offboarding. |
| **Department Head** | Departmental | Read/Write access for their specific department (e.g., approve leave for *their* team, view *their* cost center). |
| **Standard Employee** | Self | View own payslips, submit own leave requests, view company directory. |

## 3. Approval Hierarchies

Complex workflows use dynamic approval chains.

**Example: Purchase Order (PO) Approval Chain**

1.  **Clerk:** Creates PO ($15,000).
2.  **System:** Detects amount > $1,000. Routes to Department Manager.
3.  **Department Manager:** Approves PO.
4.  **System:** Detects amount > $10,000. Routes to Finance Controller.
5.  **Finance Controller:** Approves PO. State changes to 'APPROVED'.

## 4. Implementation Strategy (Backend)

Permissions are evaluated at the API gateway / middleware level using JWT claims.

*   JWT Payload includes: `user_id`, `role_ids`, `department_ids`.
*   Middleware checks if the user's roles possess the required permission (e.g., `invoice:approve`) for the requested resource context.