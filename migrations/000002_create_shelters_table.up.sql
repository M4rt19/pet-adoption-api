CREATE TABLE IF NOT EXISTS shelters (
                                        id SERIAL PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        address TEXT,
                                        phone TEXT,
                                        owner_user_id INT NOT NULL,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_shelters_owner_user
    FOREIGN KEY (owner_user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_shelters_owner_user_id ON shelters (owner_user_id);
