ALTER TABLE employee_docs RENAME COLUMN doc_type TO doc_name;
ALTER TABLE employee_docs ADD COLUMN issued_at BIGINT NOT NULL DEFAULT 0;
ALTER TABLE employee_docs DROP COLUMN IF EXISTS doc_number;
