ALTER TABLE visits
ADD COLUMN IF NOT EXISTS date VARCHAR(20);

ALTER TABLE visits
ADD COLUMN IF NOT EXISTS updated_at BIGINT NOT NULL DEFAULT 0;

ALTER TABLE visits
ALTER COLUMN created_at TYPE BIGINT
USING created_at::BIGINT;

UPDATE visits
SET date = TO_CHAR(TO_TIMESTAMP(created_at / 1000.0), 'YYYY-MM-DD')
WHERE date IS NULL OR date = '';

UPDATE visits
SET updated_at = created_at
WHERE updated_at IS NULL;

ALTER TABLE visits
ALTER COLUMN date SET NOT NULL;

CREATE TABLE IF NOT EXISTS visit_details (
    id VARCHAR(36) NOT NULL,
    visit_id VARCHAR(36) NOT NULL,
    visit_type VARCHAR(10) NOT NULL,
    visit_at VARCHAR(20) NOT NULL,
    date_visit VARCHAR(20) NOT NULL,
    file_url TEXT,
    location TEXT,
    address TEXT,
    note TEXT,
    created_at BIGINT NOT NULL DEFAULT 0,
    updated_at BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT fk_visit_details_visit_id
        FOREIGN KEY (visit_id) REFERENCES visits (id) ON DELETE CASCADE
);

INSERT INTO visit_details (
    id,
    visit_id,
    visit_type,
    visit_at,
    date_visit,
    file_url,
    location,
    address,
    note,
    created_at,
    updated_at
)
SELECT
    SUBSTRING(MD5(visits.id || '-legacy-detail') FROM 1 FOR 36),
    visits.id,
    visits.visit_type,
    TO_CHAR(TO_TIMESTAMP(visits.created_at / 1000.0), 'HH24:MI:SS'),
    TO_CHAR(TO_TIMESTAMP(visits.created_at / 1000.0), 'YYYY-MM-DD'),
    visits.file_url,
    CASE
        WHEN visits.latitude IS NOT NULL AND visits.longitude IS NOT NULL
            THEN visits.latitude || ', ' || visits.longitude
        ELSE NULL
    END,
    visits.address,
    visits.note,
    visits.created_at,
    visits.created_at
FROM visits
WHERE visits.visit_type IS NOT NULL
AND NOT EXISTS (
    SELECT 1
    FROM visit_details
    WHERE visit_details.visit_id = visits.id
);

CREATE INDEX IF NOT EXISTS idx_visit_details_visit_id ON visit_details (visit_id);
CREATE INDEX IF NOT EXISTS idx_visit_details_visit_type ON visit_details (visit_type);

ALTER TABLE visits
DROP COLUMN IF EXISTS visit_type,
DROP COLUMN IF EXISTS note,
DROP COLUMN IF EXISTS latitude,
DROP COLUMN IF EXISTS longitude,
DROP COLUMN IF EXISTS address,
DROP COLUMN IF EXISTS file_name,
DROP COLUMN IF EXISTS mime_type,
DROP COLUMN IF EXISTS file_size,
DROP COLUMN IF EXISTS file_url;
