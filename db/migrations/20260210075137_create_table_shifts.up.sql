CREATE TABLE shifts 
(
    id VARCHAR(100) NOT NULL,
    company_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    late_tolerance INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(id),

    CONSTRAINT fk_shift_company_id
          FOREIGN KEY (company_id)
          REFERENCES companies(id)
          ON DELETE CASCADE
);

CREATE TABLE employee_shifts 
(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    employee_id VARCHAR(36) NOT NULL,
    shift_id VARCHAR(100) NOT NULL,
--     assigned_date DATE NOT NULL,
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
