-- roles: entity has CompanyID but table was created without it
ALTER TABLE roles
    ADD COLUMN IF NOT EXISTS company_id VARCHAR(100) NOT NULL DEFAULT '';

-- time_off_requests: entity has CompanyID but table was created without it
ALTER TABLE time_off_requests
    ADD COLUMN IF NOT EXISTS company_id VARCHAR(36) NOT NULL DEFAULT '';

-- time_off_balances: entity has CompanyID but table was created without it
ALTER TABLE time_off_balances
    ADD COLUMN IF NOT EXISTS company_id VARCHAR(36) NOT NULL DEFAULT '';

-- employee_docs: entity has CompanyID but table was created without it
ALTER TABLE employee_docs
    ADD COLUMN IF NOT EXISTS company_id VARCHAR(100) NOT NULL DEFAULT '';
