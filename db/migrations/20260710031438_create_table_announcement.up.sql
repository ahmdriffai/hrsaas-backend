CREATE TABLE announcements (
    id VARCHAR(100) NOT NULL,
    company_id VARCHAR(100) NOT NULL,
    employee_id VARCHAR(100) NOT NULL,
    title VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (id),

    CONSTRAINT fk_announcements_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_announcements_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
);