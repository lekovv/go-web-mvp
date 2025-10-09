CREATE TABLE appointment_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO appointment_statuses (name) VALUES
    ('scheduled'),
    ('confirmed'),
    ('completed'),
    ('cancelled'),
    ('no_show');