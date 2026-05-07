CREATE TABLE state_admin (
    id SERIAL PRIMARY KEY,

    profile_id UUID NOT NULL,      -- from profiles table
    assigned_state INT NOT NULL,   -- from state table

    CONSTRAINT fk_profile
        FOREIGN KEY (profile_id)
        REFERENCES profiles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_assigned_state
        FOREIGN KEY (assigned_state)
        REFERENCES state(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);