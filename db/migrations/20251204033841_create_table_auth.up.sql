-- USER TABLE
CREATE TABLE users (
  id VARCHAR(36) NOT NULL,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(255) NOT NULL,
  email_verified BOOLEAN DEFAULT FALSE,
  image VARCHAR(255) NULL ,
  role VARCHAR(50) DEFAULT 'USER',
  company_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id),

  CONSTRAINT fk_user_company_id 
    FOREIGN KEY (company_id) REFERENCES companies (id)
);

-- SESSION TABLE
CREATE TABLE sessions (
  id VARCHAR(36) NOT NULL,
  token TEXT NOT NULL,
  ip_address VARCHAR(50),
  user_agent TEXT,
  user_id VARCHAR(36) NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  CONSTRAINT fk_session_user_id 
    FOREIGN KEY (user_id) REFERENCES users (id)
);
