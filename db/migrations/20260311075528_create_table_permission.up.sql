CREATE TABLE time_off_types
(
    id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(10) NOT NULL,
    is_quota_based BOOLEAN NOT NULL DEFAULT FALSE,
    default_quota_days INT NOT NULL DEFAULT 0,
    PRIMARY KEY(id)
);

CREATE TABLE time_off_requests(
    id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    time_off_type_id VARCHAR(50) NOT NULL,
    start_date BIGINT NOT NULL,
    end_date BIGINT NOT NULL,
    request_reason TEXT,
    request_status VARCHAR(50) NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY(id),

    CONSTRAINT fk_time_off_request_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_time_off_request_time_off_type_id
        FOREIGN KEY (time_off_type_id)
        REFERENCES time_off_types(id)
        ON DELETE CASCADE
);

CREATE TABLE time_off_balances
(
    id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    time_off_type_id VARCHAR(36) NOT NULL,
    period_year INT NOT NULL,
    entitled_days INT NOT NULL,
    used_days INT NOT NULL,
    remaining_days INT NOT NULL,
    PRIMARY KEY(id),

    CONSTRAINT fk_time_off_balance_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_time_off_balance_time_off_type_id
        FOREIGN KEY (time_off_type_id)
        REFERENCES time_off_types(id)
        ON DELETE CASCADE
);

CREATE TABLE time_off_attachments
(
    id VARCHAR(36) NOT NULL,
    time_off_request_id VARCHAR(36) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size INT NOT NULL,
    file_url TEXT NOT NULL,
    PRIMARY KEY(id),

    CONSTRAINT fk_time_off_attachment_request_id
        FOREIGN KEY (time_off_request_id)
        REFERENCES time_off_requests(id)
        ON DELETE CASCADE
);

CREATE TABLE time_off_approvals
(
    id VARCHAR(36) NOT NULL,
    time_off_request_id VARCHAR(36) NOT NULL,
    approver_id VARCHAR(36) NOT NULL,
    approval_status VARCHAR(50) NOT NULL,
    action_reason TEXT,
    action_at BIGINT NULL,
    PRIMARY KEY(id),

    CONSTRAINT fk_time_off_approval_request_id
        FOREIGN KEY (time_off_request_id)
        REFERENCES time_off_requests(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_time_off_approval_approver_id
        FOREIGN KEY (approver_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
);
