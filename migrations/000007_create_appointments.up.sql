CREATE TABLE appointments (
    id UUID PRIMARY KEY,
    doctor_id UUID NOT NULL,
    patient_id UUID NOT NULL,
    appointment_date TIMESTAMP NOT NULL,
    appointment_status_id UUID NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    CONSTRAINT fk_appointments_doctor_id FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointments_patient_id FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointments_status_id FOREIGN KEY (appointment_status_id) REFERENCES appointment_statuses(id) ON DELETE RESTRICT
);

CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_date ON appointments(appointment_date);
CREATE INDEX idx_appointments_status_id ON appointments(appointment_status_id);