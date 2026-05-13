ALTER TABLE roles
    DROP COLUMN IF EXISTS company_id;

ALTER TABLE time_off_requests
    DROP COLUMN IF EXISTS company_id;

ALTER TABLE time_off_balances
    DROP COLUMN IF EXISTS company_id;

ALTER TABLE employee_docs
    DROP COLUMN IF EXISTS company_id;
