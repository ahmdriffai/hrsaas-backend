CREATE TABLE attendances 
(
  id VARCHAR(100) NOT NULL,
  company_id VARCHAR(36) NOT NULL,
  employee_id VARCHAR(36) NOT NULL,
  date DATE NOT NULL,
  check_in_time TIMESTAMPTZ,
  check_out_time TIMESTAMPTZ,
  total_work_minutes INT,
  total_break_minutes INT,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY(id),

  CONSTRAINT fk_attendance_company_id
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
        ON DELETE CASCADE,

  CONSTRAINT fk_attendance_employee_id
        FOREIGN KEY (employee_id)
        REFERENCES employees(id)
        ON DELETE CASCADE
);

CREATE TABLE attendance_logs
(
  id VARCHAR(100) NOT NULL,
  attendance_id VARCHAR(100) NOT NULL,
  type VARCHAR(50) NOT NULL,
  time TIMESTAMPTZ NOT NULL,
  lat DOUBLE PRECISION,
  lng DOUBLE PRECISION,
  location_distance DOUBLE PRECISION,
  is_location_verified BOOLEAN NOT NULL DEFAULT FALSE,
  is_face_verified BOOLEAN NOT NULL DEFAULT FALSE,
  face_confidence DOUBLE PRECISION,
  face_image_url TEXT,
  device_info TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY(id),

  CONSTRAINT fk_attendance_log_attendance_id
        FOREIGN KEY (attendance_id)
        REFERENCES attendances(id)
        ON DELETE CASCADE
);