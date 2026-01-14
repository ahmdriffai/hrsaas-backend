CREATE TABLE companies 
(
  id VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  logo_url VARCHAR(255),
  bussiness_field VARCHAR(20),
  address VARCHAR(255),
  province VARCHAR(50),
  city VARCHAR(50),
  district VARCHAR(50),
  village VARCHAR(50),
  zip_code VARCHAR(50),
  phone_number VARCHAR(20),
  fax_number VARCHAR(20),
  email VARCHAR(100),
  website VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY(id)
);

-- DIVISIONS AND POSITIONS TABLES ARE ASSUMED TO BE CREATED BEFORE THIS MIGRATION

CREATE TABLE divisions 
(
  id VARCHAR(100) NOT NULL,
  company_id VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  CONSTRAINT fk_division_company_id 
    FOREIGN KEY (company_id) REFERENCES companies (id)
);

CREATE TABLE positions (
  id VARCHAR(36),
  company_id VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT NULL,
  parent_id VARCHAR(50) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  CONSTRAINT fk_parent
      FOREIGN KEY (parent_id) REFERENCES positions(id)
      ON DELETE SET NULL
)