ALTER TABLE employee_educations
    ALTER COLUMN graduation_year TYPE BIGINT USING EXTRACT(EPOCH FROM graduation_year)::BIGINT * 1000;