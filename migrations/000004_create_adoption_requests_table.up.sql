CREATE TABLE IF NOT EXISTS adoption_requests (
                                                 id SERIAL PRIMARY KEY,
                                                 user_id INT NOT NULL,
                                                 pet_id INT NOT NULL,
                                                 status TEXT NOT NULL DEFAULT 'pending'
                                                 CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled', 'expired')),
    message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_adoption_requests_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE,

    CONSTRAINT fk_adoption_requests_pet
    FOREIGN KEY (pet_id)
    REFERENCES pets (id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_adoption_requests_user_id ON adoption_requests (user_id);
CREATE INDEX IF NOT EXISTS idx_adoption_requests_pet_id ON adoption_requests (pet_id);
CREATE INDEX IF NOT EXISTS idx_adoption_requests_status ON adoption_requests (status);
