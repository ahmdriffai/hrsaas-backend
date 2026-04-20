ALTER TABLE positions
ADD COLUMN is_approver BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE time_off_approvals
ADD COLUMN is_required BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX idx_positions_is_approver ON positions(is_approver);
