CREATE TABLE IF NOT EXISTS pets (
                                    id SERIAL PRIMARY KEY,
                                    shelter_id INT NOT NULL,
                                    name TEXT NOT NULL,
                                    species TEXT NOT NULL,         -- e.g. dog, cat
                                    breed TEXT,
                                    age INT,
                                    description TEXT,
                                    status TEXT NOT NULL DEFAULT 'available'
                                    CHECK (status IN ('available', 'reserved', 'adopted')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_pets_shelter
    FOREIGN KEY (shelter_id)
    REFERENCES shelters (id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_pets_shelter_id ON pets (shelter_id);
CREATE INDEX IF NOT EXISTS idx_pets_status ON pets (status);
