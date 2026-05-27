CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

--------------------------------------------------------------------------------------------------------------
-- ROLES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE roles (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    code CITEXT UNIQUE NOT NULL
        CHECK (length(trim(code)) > 0),

    display_name VARCHAR(50) NOT NULL
        CHECK (length(trim(display_name)) > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------------------------------------------------------------------
-- STATES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE states (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    name CITEXT UNIQUE NOT NULL
        CHECK (length(trim(name)) > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICTS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE districts (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    name CITEXT NOT NULL
        CHECK (length(trim(name)) > 0),

    state_id SMALLINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_districts_state_id
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT unique_district_per_state
        UNIQUE(name, state_id)
);

CREATE INDEX idx_districts_state_id
ON districts(state_id);

--------------------------------------------------------------------------------------------------------------
-- PINCODES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE pincodes (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    code CHAR(6) UNIQUE NOT NULL
        CHECK (code ~ '^[0-9]{6}$'),

    district_id INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pincodes_district_id
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX idx_pincodes_district_id
ON pincodes(district_id);


--------------------------------------------------------------------------------------------------------------
-- DISCIPLINES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE disciplines (

    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    code CITEXT UNIQUE NOT NULL
        CHECK (length(trim(code)) > 0),

    display_name VARCHAR(50) NOT NULL
        CHECK (length(trim(display_name)) > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------------------------------------------------------------------
-- USERS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name VARCHAR(50) NOT NULL
        CHECK (length(trim(first_name)) > 0),

    last_name VARCHAR(50)
        CHECK (
            last_name IS NULL
            OR length(trim(last_name)) > 0
        ),

    email CITEXT UNIQUE NOT NULL
        CHECK (length(trim(email)) > 0),

    password_hash TEXT NOT NULL,

    contact_number VARCHAR(20) NOT NULL
        CHECK (length(trim(contact_number)) > 0),

    role_id SMALLINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_users_role_id
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX idx_users_role_id
ON users(role_id);

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
-- USERS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- ACADEMIES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE academies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name CITEXT NOT NULL
        CHECK (length(trim(name)) > 0),

    district_id INTEGER NOT NULL,

    address TEXT NOT NULL
        CHECK (length(trim(address)) > 0),

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_academies_district_id
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMIES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_academies_district_id
ON academies(district_id);

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

    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    email CITEXT NOT NULL
        CHECK (length(trim(email)) > 0),

    role_id SMALLINT NOT NULL,

    invited_by UUID NOT NULL,

    token_hash TEXT UNIQUE NOT NULL,

    scope_type VARCHAR(50),

    scope_id TEXT,

    expires_at TIMESTAMPTZ NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (
            status IN (
                'pending',
                'accepted',
                'expired',
                'revoked'
            )
        ),

    accepted_at TIMESTAMPTZ,

    used_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_invitations_role_id
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT fk_invitations_invited_by
        FOREIGN KEY (invited_by)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_invitations_used_by
        FOREIGN KEY (used_by)
        REFERENCES users(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- INVITATIONS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_invitations_email
ON invitations(email);

CREATE INDEX idx_invitations_role_id
ON invitations(role_id);

CREATE INDEX idx_invitations_invited_by
ON invitations(invited_by);

CREATE INDEX idx_invitations_used_by
ON invitations(used_by);

CREATE INDEX idx_invitations_status
ON invitations(status);

--------------------------------------------------------------------------------------------------------------
-- INVITATIONS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER invitations_set_updated_at
BEFORE UPDATE ON invitations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE state_admins (

    user_id UUID PRIMARY KEY,

    state_id SMALLINT NOT NULL,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,

    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_state_admins_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_state_admins_state_id
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_state_admins_state_id
ON state_admins(state_id);

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER state_admins_set_updated_at
BEFORE UPDATE ON state_admins
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- DISTRICT ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE district_admins (

    user_id UUID PRIMARY KEY,

    district_id INTEGER NOT NULL,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,
    
    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_district_admins_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_district_admins_district_id
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT ADMINS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_district_admins_district_id
ON district_admins(district_id);

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

    user_id UUID PRIMARY KEY,

    district_id INTEGER NOT NULL,

    coach_code VARCHAR(20)
        CHECK (
            coach_code IS NULL
            OR length(trim(coach_code)) > 0
        ),

    coaching_certificate_proof TEXT,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,

    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_district_coach_code
        UNIQUE (district_id, coach_code),

    CONSTRAINT fk_district_coaches_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_district_coaches_district_id
        FOREIGN KEY (district_id)
        REFERENCES districts(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACHES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_district_coaches_district_id
ON district_coaches(district_id);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACHES UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER district_coaches_set_updated_at
BEFORE UPDATE ON district_coaches
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACH DISCIPLINES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE district_coach_disciplines (

    coach_user_id UUID NOT NULL,

    discipline_id SMALLINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (
        coach_user_id,
        discipline_id
    ),

    CONSTRAINT fk_district_coach_disciplines_coach_user_id
        FOREIGN KEY (coach_user_id)
        REFERENCES district_coaches(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_district_coach_disciplines_discipline_id
        FOREIGN KEY (discipline_id)
        REFERENCES disciplines(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACH DISCIPLINES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_district_coach_disciplines_discipline_id
ON district_coach_disciplines(discipline_id);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY ADMINS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE academy_admins (

    user_id UUID PRIMARY KEY,

    academy_id UUID NOT NULL,

    gstin VARCHAR(15),

    registration_proof TEXT,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,

    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_academy_admins_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_academy_admins_academy_id
        FOREIGN KEY (academy_id)
        REFERENCES academies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY ADMINS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_academy_admins_academy_id
ON academy_admins(academy_id);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY ADMINS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER academy_admins_set_updated_at
BEFORE UPDATE ON academy_admins
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACHES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE academy_coaches (

    user_id UUID PRIMARY KEY,

    academy_id UUID NOT NULL,

    coach_code VARCHAR(20)
        CHECK (
            coach_code IS NULL
            OR length(trim(coach_code)) > 0
        ),

    coaching_certificate_proof TEXT,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,

    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_academy_coach_code
        UNIQUE (academy_id, coach_code),

    CONSTRAINT fk_academy_coaches_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_academy_coaches_academy_id
        FOREIGN KEY (academy_id)
        REFERENCES academies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACHES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_academy_coaches_academy_id
ON academy_coaches(academy_id);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACHES UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER academy_coaches_set_updated_at
BEFORE UPDATE ON academy_coaches
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACH DISCIPLINES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE academy_coach_disciplines (

    coach_user_id UUID NOT NULL,

    discipline_id SMALLINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (
        coach_user_id,
        discipline_id
    ),

    CONSTRAINT fk_academy_coach_disciplines_coach_user_id
        FOREIGN KEY (coach_user_id)
        REFERENCES academy_coaches(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_academy_coach_disciplines_discipline_id
        FOREIGN KEY (discipline_id)
        REFERENCES disciplines(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACH DISCIPLINES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_academy_coach_disciplines_discipline_id
ON academy_coach_disciplines(discipline_id);

--------------------------------------------------------------------------------------------------------------
-- GENDER TYPE
--------------------------------------------------------------------------------------------------------------

CREATE TYPE gender_type AS ENUM (
    'male',
    'female',
    'other'
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER STATUS TYPE
--------------------------------------------------------------------------------------------------------------

CREATE TYPE player_status_type AS ENUM (
    'active',
    'inactive',
    'suspended',
    'transferred'
);

--------------------------------------------------------------------------------------------------------------
-- PLAYERS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE players (

    user_id UUID PRIMARY KEY,

    academy_id UUID NOT NULL,

    current_coach_user_id UUID,

    status player_status_type NOT NULL DEFAULT 'active',

    joined_at DATE NOT NULL DEFAULT CURRENT_DATE,

    dpdp_consent BOOLEAN NOT NULL DEFAULT FALSE,

    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,

    registered_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_players_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_players_academy_id
        FOREIGN KEY (academy_id)
        REFERENCES academies(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT fk_players_current_coach_user_id
        FOREIGN KEY (current_coach_user_id)
        REFERENCES academy_coaches(user_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT fk_registered_by_user_id
        FOREIGN KEY (registered_by)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYERS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_players_academy_id
ON players(academy_id);

CREATE INDEX idx_players_current_coach_user_id
ON players(current_coach_user_id);

CREATE INDEX idx_players_status
ON players(status);

CREATE INDEX idx_players_registered_by
ON players(registered_by);

--------------------------------------------------------------------------------------------------------------
-- PLAYERS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER players_set_updated_at
BEFORE UPDATE ON players
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- PLAYER PERSONAL INFO
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_personal_info (

    player_user_id UUID PRIMARY KEY,

    date_of_birth DATE NOT NULL,

    gender gender_type NOT NULL,

    nationality VARCHAR(50) NOT NULL DEFAULT 'Indian',

    place_of_birth VARCHAR(100),

    city VARCHAR(100),

    residential_address TEXT,

    pincode_id INTEGER,

    education VARCHAR(100),

    institution_name VARCHAR(150),

    occupation VARCHAR(100),

    temporary_sport_id VARCHAR(20),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_player_personal_info_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_player_personal_info_pincode_id
        FOREIGN KEY (pincode_id)
        REFERENCES pincodes(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER PERSONAL INFO INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_player_personal_info_pincode_id
ON player_personal_info(pincode_id);

--------------------------------------------------------------------------------------------------------------
-- PLAYER PERSONAL INFO UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER player_personal_info_set_updated_at
BEFORE UPDATE ON player_personal_info
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- PLAYER SPORTS PROFILE
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_sports_profile (

    player_user_id UUID PRIMARY KEY,

    unit_of_representation VARCHAR(100),

    dominant_hand VARCHAR(20)
        CHECK (
            dominant_hand IN (
                'left',
                'right',
                'ambidextrous'
            )
        ),

    height_cm DECIMAL(5,2)
        CHECK (
            height_cm IS NULL
            OR height_cm > 0
        ),

    weight_kg DECIMAL(5,2)
        CHECK (
            weight_kg IS NULL
            OR weight_kg > 0
        ),

    shoe_size VARCHAR(10),

    tracksuit_size VARCHAR(10),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_player_sports_profile_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER SPORTS PROFILE UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER player_sports_profile_set_updated_at
BEFORE UPDATE ON player_sports_profile
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- PLAYER DISCIPLINES
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_disciplines (

    player_user_id UUID NOT NULL,

    discipline_id SMALLINT NOT NULL,

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (
        player_user_id,
        discipline_id
    ),

    CONSTRAINT fk_player_disciplines_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_player_disciplines_discipline_id
        FOREIGN KEY (discipline_id)
        REFERENCES disciplines(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER DISCIPLINES INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_player_disciplines_discipline_id
ON player_disciplines(discipline_id);

CREATE UNIQUE INDEX unique_primary_player_discipline
ON player_disciplines(player_user_id)
WHERE is_primary = TRUE;

--------------------------------------------------------------------------------------------------------------
-- PLAYER GUARDIANS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_guardians (

    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    player_user_id UUID NOT NULL,

    full_name VARCHAR(100) NOT NULL
        CHECK (
            length(trim(full_name)) BETWEEN 2 AND 100
        ),

    relationship VARCHAR(50),

    contact_number VARCHAR(20) NOT NULL
        CHECK (
            length(trim(contact_number)) > 0
        ),

    alternative_contact VARCHAR(20),

    parental_consent BOOLEAN NOT NULL DEFAULT FALSE,

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_player_guardians_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER GUARDIANS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_player_guardians_player_user_id
ON player_guardians(player_user_id);

CREATE UNIQUE INDEX unique_primary_player_guardian
ON player_guardians(player_user_id)
WHERE is_primary = TRUE;

--------------------------------------------------------------------------------------------------------------
-- PLAYER PASSPORTS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_passports (

    player_user_id UUID PRIMARY KEY,

    passport_number VARCHAR(20) UNIQUE,

    passport_issue_date DATE,

    passport_expiry_date DATE,

    passport_issuing_authority VARCHAR(100),

    passport_place_of_issue VARCHAR(100),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_player_passports_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER PASSPORTS UPDATED_AT TRIGGER
--------------------------------------------------------------------------------------------------------------

CREATE TRIGGER player_passports_set_updated_at
BEFORE UPDATE ON player_passports
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

--------------------------------------------------------------------------------------------------------------
-- PLAYER COACH ASSIGNMENTS
--------------------------------------------------------------------------------------------------------------

CREATE TABLE player_coach_assignments (

    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    player_user_id UUID NOT NULL,

    coach_user_id UUID NOT NULL,

    assigned_from DATE NOT NULL,

    assigned_until DATE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_player_coach_assignments_player_user_id
        FOREIGN KEY (player_user_id)
        REFERENCES players(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_player_coach_assignments_coach_user_id
        FOREIGN KEY (coach_user_id)
        REFERENCES academy_coaches(user_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT check_assignment_dates
        CHECK (
            assigned_until IS NULL
            OR assigned_until >= assigned_from
        )
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER COACH ASSIGNMENTS INDEXES
--------------------------------------------------------------------------------------------------------------

CREATE INDEX idx_player_coach_assignments_player_user_id
ON player_coach_assignments(player_user_id);

CREATE INDEX idx_player_coach_assignments_coach_user_id
ON player_coach_assignments(coach_user_id);

