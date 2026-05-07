ALTER TABLE visits
ADD COLUMN IF NOT EXISTS visit_type VARCHAR(10),
ADD COLUMN IF NOT EXISTS note TEXT,
ADD COLUMN IF NOT EXISTS latitude VARCHAR(50),
ADD COLUMN IF NOT EXISTS longitude VARCHAR(50),
ADD COLUMN IF NOT EXISTS address TEXT,
ADD COLUMN IF NOT EXISTS file_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
ADD COLUMN IF NOT EXISTS file_size INT,
ADD COLUMN IF NOT EXISTS file_url TEXT;

UPDATE visits
SET
    visit_type = details.visit_type,
    note = details.note,
    address = details.address,
    file_url = details.attachment_url
FROM (
    SELECT DISTINCT ON (visit_id)
        visit_id,
        visit_type,
        note,
        address,
        attachment_url
    FROM visit_details
    ORDER BY visit_id, created_at ASC
) AS details
WHERE visits.id = details.visit_id;

DROP INDEX IF EXISTS idx_visit_details_visit_type;
DROP INDEX IF EXISTS idx_visit_details_visit_id;
DROP TABLE IF EXISTS visit_details;

ALTER TABLE visits
DROP COLUMN IF EXISTS date,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE visits
ALTER COLUMN visit_type SET NOT NULL;
