ALTER TABLE employee_educations
    ADD COLUMN company_id  VARCHAR(36) NULL,
    ADD COLUMN start_year  INT         NULL,
    ADD COLUMN end_year    INT         NULL,
    ADD COLUMN created_at  BIGINT      NULL,
    ADD COLUMN updated_at  BIGINT      NULL;

ALTER TABLE employee_educations
    ADD CONSTRAINT fk_education_company_id
        FOREIGN KEY (company_id) REFERENCES companies (id);
