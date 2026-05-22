CREATE TABLE holidays
(
    id VARCHAR(100) NOT NULL,
    company_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_national BOOLEAN,
    date BIGINT NOT NULL,
    end_date BIGINT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,

    PRIMARY KEY(id)
);