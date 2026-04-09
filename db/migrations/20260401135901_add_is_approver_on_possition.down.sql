DROP INDEX IF EXISTS idx_positions_is_approver;

ALTER TABLE time_off_approvals
DROP COLUMN IF EXISTS is_required;

ALTER TABLE positions
DROP COLUMN IF EXISTS is_approver;
