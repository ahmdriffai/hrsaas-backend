CREATE TABLE employee_trainings (
    id              VARCHAR(100) NOT NULL,
    company_id      VARCHAR(100) NOT NULL,
    employee_id     VARCHAR(100) NOT NULL,
    training_name   VARCHAR(255) NOT NULL,
    organizer       VARCHAR(255) NOT NULL,
    start_date      BIGINT       NOT NULL,
    end_date        BIGINT,
    certificate_url VARCHAR(255),
    created_at      BIGINT       NOT NULL,
    updated_at      BIGINT       NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT fk_employee_trainings_employee
        FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE
);
