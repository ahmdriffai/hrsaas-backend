ALTER TABLE employee_docs RENAME COLUMN doc_name TO doc_type;
ALTER TABLE employee_docs DROP COLUMN IF EXISTS issued_at;
ALTER TABLE employee_docs ADD COLUMN doc_number VARCHAR(100) NOT NULL DEFAULT '';
