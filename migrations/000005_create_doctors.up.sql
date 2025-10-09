CREATE TABLE doctors (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    specialization_id UUID NOT NULL,
    bio TEXT,
    experience_years INTEGER,
    price INTEGER NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    CONSTRAINT fk_doctors_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_doctors_specialization_id FOREIGN KEY (specialization_id) REFERENCES specializations(id) ON DELETE RESTRICT
);

CREATE INDEX idx_doctors_user_id ON doctors(user_id);
CREATE INDEX idx_doctors_specialization_id ON doctors(specialization_id);
CREATE INDEX idx_doctors_price ON doctors(price);