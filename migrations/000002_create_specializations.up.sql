CREATE TABLE specializations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO specializations (name) VALUES
    ('Терапевт'),
    ('Кардиолог'),
    ('Невролог'),
    ('Дерматолог'),
    ('Офтальмолог'),
    ('Хирург'),
    ('Психиатр'),
    ('Стоматолог'),
    ('Гинеколог'),
    ('Педиатр');