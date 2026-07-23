ALTER TABLE
    remidial_visit ADD COLUMN latitude VARCHAR(50);
ALTER TABLE
    remidial_visit ADD COLUMN longitude VARCHAR(50);
ALTER TABLE
    remidial_visit ADD COLUMN img_url VARCHAR(255);
ALTER TABLE remidial_visit 
    RENAME COLUMN note TO commitment;