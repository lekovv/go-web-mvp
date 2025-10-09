DROP INDEX IF EXISTS idx_appointments_status_id;
DROP INDEX IF EXISTS idx_appointments_date;
DROP INDEX IF EXISTS idx_appointments_patient_id;
DROP INDEX IF EXISTS idx_appointments_doctor_id;
DROP TABLE IF EXISTS appointments CASCADE;