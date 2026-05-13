CREATE EXTENSION IF NOT EXISTS "pgcrypto";

--------------------------------------------------------------------------------------------------------------
-- ROLES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE roles (
    id SMALLSERIAL PRIMARY KEY,

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
    id SMALLSERIAL PRIMARY KEY,

    name VARCHAR(100) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO states (name) VALUES ('Karnataka');
INSERT INTO states (name) VALUES ('Kerala');

--------------------------------------------------------------------------------------------------------------
-- DISTRICTS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE districts (
    id SERIAL PRIMARY KEY,

    name VARCHAR(50) NOT NULL,

    state_id SMALLINT NOT NULL,

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
    id SERIAL PRIMARY KEY,

    code VARCHAR(10) UNIQUE NOT NULL,

    district_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PROFILES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name VARCHAR(20) NOT NULL,

    last_name VARCHAR(20),

    email VARCHAR(50) UNIQUE NOT NULL,

    password_hash TEXT NOT NULL,

    contact_number VARCHAR(20) NOT NULL,

    role_id SMALLINT NOT NULL,

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
-- UPDATED_AT TRIGGER FUNCTION
--------------------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--------------------------------------------------------------------------------------------------------------
-- PROFILES UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER profiles_set_updated_at
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE state_admins (
    profile_id UUID PRIMARY KEY,

    state_id SMALLINT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_profile
        FOREIGN KEY (profile_id)
        REFERENCES profiles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_state
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE district_admins (
    profile_id UUID PRIMARY KEY,

    state_id SMALLINT NOT NULL,

    district_id INT NOT NULL,

    dpdp_consent BOOLEAN NOT NULL
        CHECK (dpdp_consent = TRUE),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_profile
        FOREIGN KEY (profile_id)
        REFERENCES profiles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_state
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT ADMINS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER district_admins_set_updated_at
BEFORE UPDATE ON district_admins
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACHES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE district_coaches (

    profile_id UUID NOT PRIMARY KEY,
    district_id INT NOT NULL,

    coach_code VARCHAR(20) UNIQUE,
    coaching_certificate_proof TEXT,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_profiles
        FOREIGN KEY (profile_id)
        REFERENCES profiles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_districts
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);


--------------------------------------------------------------------------------------------------------------
-- ACADEMIES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE academies (
    id SERIAL PRIMARY KEY,

    name VARCHAR(50) NOT NULL,

    district_id INT NOT NULL,

    address TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMIES UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER academies_set_updated_at
BEFORE UPDATE ON academies
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- INVITATIONS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE invitations (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(50) NOT NULL,

    role_id SMALLINT NOT NULL,

    invited_by UUID NOT NULL,

    token TEXT NOT NULL UNIQUE,

    state_id SMALLINT,
    district_id INT,
    academy_id INT,

    expires_at TIMESTAMP NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted')),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_invited_by
        FOREIGN KEY (invited_by)
        REFERENCES profiles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_state
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_academy
        FOREIGN KEY (academy_id)
        REFERENCES academies(id)
        ON DELETE RESTRICT
);

