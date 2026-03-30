CREATE TABLE visits (
    id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    company_id VARCHAR(36) NOT NULL,
    visit_type VARCHAR(10) NOT NULL,
    note TEXT,
    latitude VARCHAR(50),
    longitude VARCHAR(50),
    address TEXT,
    file_name VARCHAR(255),
    mime_type VARCHAR(100),
    file_size INT,
    file_url TEXT,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (id),

    CONSTRAINT fk_visits_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_visits_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_visits_employee_id ON visits(employee_id);
CREATE INDEX idx_visits_company_id ON visits(company_id);
CREATE INDEX idx_visits_created_at ON visits(created_at);