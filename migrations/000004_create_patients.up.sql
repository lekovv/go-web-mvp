CREATE TABLE patients (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    birth_date DATE NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    CONSTRAINT fk_patients_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_patients_user_id ON patients(user_id);
CREATE INDEX idx_patients_birth_date ON patients(birth_date);