CREATE EXTENSION IF NOT EXISTS "pgcrypto";

--------------------------------------------------------------------------------------------------------------
-- ROLES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(20) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name) VALUES ('super_admin');
INSERT INTO roles (name) VALUES ('state_admin');
INSERT INTO roles (name) VALUES ('district_admin');
INSERT INTO roles (name) VALUES ('academic_admin');
INSERT INTO roles (name) VALUES ('coach');
INSERT INTO roles (name) VALUES ('player');

--------------------------------------------------------------------------------------------------------------
-- STATES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO states (name) VALUES ('Karnataka');
INSERT INTO states (name) VALUES ('Kerala');

--------------------------------------------------------------------------------------------------------------
-- DISTRICTS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(50) NOT NULL,

    state_id UUID NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_state
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT unique_district_per_state
        UNIQUE(name, state_id)
);

INSERT INTO districts (name, state_id)
VALUES
(
    'Bengaluru Urban',
    (SELECT id FROM states WHERE name = 'Karnataka')
),
(
    'Chitradurga',
    (SELECT id FROM states WHERE name = 'Karnataka')
),
(
    'Davangere',
    (SELECT id FROM states WHERE name = 'Karnataka')
),
(
    'Mysuru',
    (SELECT id FROM states WHERE name = 'Karnataka')
);

--------------------------------------------------------------------------------------------------------------
-- PINCODES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE pincodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code VARCHAR(10) UNIQUE NOT NULL,

    district_id UUID NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PROFILES
--------------------------------------------------------------------------------------------------------------

-- Normalize email in application layer

CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name VARCHAR(20) NOT NULL,

    last_name VARCHAR(20),

    email VARCHAR(50) UNIQUE NOT NULL,

    password_hash TEXT NOT NULL,

    contact_number VARCHAR(20) NOT NULL,

    role_id UUID NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

INSERT INTO profiles (
    first_name,
    last_name,
    email,
    password_hash,
    contact_number,
    role_id
)
VALUES (
    'Bhaskar',
    'vb',
    'bhaskarvb@gmail.com',
    '$2a$10$NXPrHEt7doY0e2Y1FgyoX.YAwP81O3UGN8Kn1/8MPwPbWxu.sSl1K',
    '1234567890',
    (
        SELECT id
        FROM roles
        WHERE name = 'super_admin'
    )
);

--------------------------------------------------------------------------------------------------------------
-- UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER profiles_set_updated_at
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE state_admins (
    profile_id UUID PRIMARY KEY,

    assigned_state_id UUID NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_profile
        FOREIGN KEY (profile_id)
        REFERENCES profiles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_assigned_state
        FOREIGN KEY (assigned_state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- Invitations
--------------------------------------------------------------------------------------------------------------

CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(50) NOT NULL,

    role_id UUID NOT NULL,

    invited_by UUID NOT NULL,

    token TEXT NOT NULL UNIQUE,

    assigned_state_id UUID,

    expires_at TIMESTAMP NOT NULL,

    used BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_invited_by
        FOREIGN KEY (invited_by)
        REFERENCES profiles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assigned_state
        FOREIGN KEY (assigned_state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
);