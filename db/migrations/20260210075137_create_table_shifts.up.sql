CREATE TABLE shifts 
(
    id VARCHAR(100) NOT NULL,
    company_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    late_tolerance INT NOT NULL,
    schedule_type VARCHAR(20) NOT NULL DEFAULT 'fixed', -- fixed | flexible
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(id),

    CONSTRAINT fk_shift_company_id
          FOREIGN KEY (company_id)
          REFERENCES companies(id)
          ON DELETE CASCADE
);

CREATE TABLE shift_days (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    shift_id VARCHAR(100) NOT NULL,
    weekday SMALLINT NOT NULL CHECK (weekday BETWEEN 1 AND 7), -- 1 Senin ... 7 Minggu
    day_type VARCHAR(20) NOT NULL DEFAULT 'workday', -- workday | offday
    check_in TIME NULL,
    check_out TIME NULL,
    break_start TIME NULL,
    break_end TIME NULL,
    max_break_minutes INT NOT NULL DEFAULT 60,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_shift_days_shift_id
    FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE CASCADE,
    CONSTRAINT uq_shift_weekday UNIQUE (shift_id, weekday)
);

CREATE TABLE employee_shifts 
(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    employee_id VARCHAR(36) NOT NULL,
    shift_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(id),

    CONSTRAINT fk_employee_shift_employee_id
          FOREIGN KEY (employee_id)
          REFERENCES employees(id)
          ON DELETE CASCADE,

    CONSTRAINT fk_employee_shift_shift_id
          FOREIGN KEY (shift_id)
          REFERENCES shifts(id)
          ON DELETE CASCADE
);
