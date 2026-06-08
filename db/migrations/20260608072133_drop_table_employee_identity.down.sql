CREATE TABLE IF NOT EXISTS employee_identities (
  id VARCHAR(36) NOT NULL,
  employee_id VARCHAR(36) NOT NULL,
  identity_type VARCHAR(20) NOT NULL,
  identity_number VARCHAR(50) NOT NULL,
  address TEXT NOT NULL,
  city VARCHAR(100) NOT NULL,
  postal_code VARCHAR(20) NOT NULL,
  domicile_address TEXT NULL,
  domicily_as_ktp BOOLEAN DEFAULT FALSE,
  is_default BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (id),
  CONSTRAINT fk_identity_employee_id
    FOREIGN KEY (employee_id) REFERENCES employees (id)
);

ALTER TABLE employees DROP COLUMN IF EXISTS identity_number;
ALTER TABLE employees DROP COLUMN IF EXISTS address;
ALTER TABLE employees DROP COLUMN IF EXISTS city;
