ALTER TABLE employee_educations
    DROP CONSTRAINT IF EXISTS fk_education_company_id,
    DROP COLUMN IF EXISTS company_id,
    DROP COLUMN IF EXISTS start_year,
    DROP COLUMN IF EXISTS end_year,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
