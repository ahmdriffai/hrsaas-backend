-- TABLE EMPLOYEES
CREATE TABLE employees (
  id VARCHAR(36) NOT NULL,
  company_id VARCHAR(36) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  employee_number VARCHAR(50) NOT NULL,
  fullname VARCHAR(100) NOT NULL,
  gender VARCHAR(10) NOT NULL,
  birth_place VARCHAR(100) NOT NULL,
  birth_date DATE NOT NULL,
  blood_type VARCHAR(5) NULL,
  marital_status VARCHAR(50) NOT NULL,
  religion VARCHAR(50) NOT NULL,
  phone VARCHAR(20) NOT NULL,
  timezone VARCHAR(50) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_employee_company_id 
    FOREIGN KEY (company_id) REFERENCES companies (id),
  CONSTRAINT fk_employee_user_id 
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- IDENTITY TABLE FOR EMPLOYEES
CREATE TABLE employee_identities (
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

-- EDUCATION TABLE FOR EMPLOYEES
CREATE TABLE employee_educations (
  id VARCHAR(36) NOT NULL,
  employee_id VARCHAR(36) NOT NULL,
  education_level VARCHAR(20) NOT NULL,
  institution_name VARCHAR(100) NOT NULL,
  major VARCHAR(100) NOT NULL,
  graduation_year DATE NOT NULL,
  gpa DECIMAL(3,2) NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_education_employee_id 
    FOREIGN KEY (employee_id) REFERENCES employees (id)
);

-- CONTRACT TABLE FOR EMPLOYEES
CREATE TABLE employee_contracts (
  id VARCHAR(36) NOT NULL,
  employee_id VARCHAR(36) NOT NULL,
  contract_type VARCHAR(20) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  division_id VARCHAR(100) NOT NULL,
  position_id VARCHAR(100) NOT NULL,
  salary DECIMAL(15,2) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_contract_employee_id 
    FOREIGN KEY (employee_id) REFERENCES employees (id),

  CONSTRAINT fk_contract_division_id 
    FOREIGN KEY (division_id) REFERENCES divisions (id),

  CONSTRAINT fk_contract_position_id 
    FOREIGN KEY (position_id) REFERENCES positions (id)
);