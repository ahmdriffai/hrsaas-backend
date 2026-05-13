CREATE TABLE employee_docs (
    id VARCHAR(100) NOT NULL,
    employee_id VARCHAR(100) NOT NULL,
    doc_type VARCHAR(100) NOT NULL,
    doc_number VARCHAR(100) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_employee_docs_employee
        FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE
);