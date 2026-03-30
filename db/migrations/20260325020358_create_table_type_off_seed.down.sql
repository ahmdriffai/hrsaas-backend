DELETE FROM time_off_types
WHERE id IN (
  'type-cuti-tahunan',
  'type-cuti-khusus',
  'type-sakit',
  'type-izin'
);