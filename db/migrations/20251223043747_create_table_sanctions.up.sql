--  CREATE TABLE SANCTION TYPES
CREATE TABLE sanctions (
    id VARCHAR(36),
    company_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    level INT NOT NULL, -- e.g., 1 for low, 2 for medium, 3 for high
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),

    CONSTRAINT fk_sanction_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);

-- CREATE TABLE EMPLOYEE SANCTIONS
CREATE TABLE employee_sanctions (
    id VARCHAR(36) ,
    employee_id VARCHAR(36) NOT NULL,
    sanction_id VARCHAR(36) NOT NULL,
    company_id VARCHAR(36) NOT NULL,
    document_url VARCHAR(255) NOT NULL,
    reason TEXT,
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'active', -- e.g., active, lifted, expired
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),

    CONSTRAINT fk_employee_sanction_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_employee_sanction_sanction_id
        FOREIGN KEY (sanction_id)
        REFERENCES sanctions(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_employee_sanction_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);