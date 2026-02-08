CREATE TABLE office_locations 
(
  id VARCHAR(100) NOT NULL,
  company_id VARCHAR(36) NOT NULL,
  name VARCHAR(100) NOT NULL,
  lat VARCHAR(255),
  lng VARCHAR(20),
  radius_meters INT,
  address TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY(id),

  CONSTRAINT fk_office_location_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);