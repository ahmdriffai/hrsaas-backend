CREATE TABLE remidial_visit (
    id VARCHAR(100) NOT NULL,
    company_id VARCHAR(100) NOT NULL,
    employee_id VARCHAR(100) NOT NULL,
    nasabah_id VARCHAR(100) NOT NULL,
    nasabah_name VARCHAR(255) NOT NULL,
    no_pjm VARCHAR(100) NOT NULL,
    loan_type VARCHAR(100) NOT NULL,
    unit VARCHAR(100) NOT NULL,
    collectibility VARCHAR(2) NOT NULL,
    loan_limit BIGINT NOT NULL,
    outstanding_balance BIGINT NOT NULL,
    overdue_principal BIGINT NOT NULL,
    overdue_interest BIGINT NOT NULL,
    overdue_total BIGINT NOT NULL,
    overdue_principal_frequency BIGINT NOT NULL,
    overdue_interest_frequency BIGINT NOT NULL,
    overdue_principal_days BIGINT NOT NULL,
    overdue_interest_days BIGINT NOT NULL,
    loan_status VARCHAR(100) NOT NULL,
    total_paid BIGINT NOT NULL,
    note TEXT,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_remidial_visit_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_remidial_visit_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
);