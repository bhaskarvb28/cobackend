--------------------------------------------------------------------------------------------------------------
-- ROLES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO roles (
    code,
    display_name
)
VALUES
    ('super_admin', 'Super Admin'),
    ('state_admin', 'State Admin'),
    ('district_admin', 'District Admin'),
    ('district_coach', 'District Coach'),
    ('academy_admin', 'Academy Admin'),
    ('academy_coach', 'Academy Coach'),
    ('player', 'Player');

--------------------------------------------------------------------------------------------------------------
-- STATES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO states (name)
VALUES
    ('Karnataka'),
    ('Kerala');

--------------------------------------------------------------------------------------------------------------
-- DISTRICTS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO districts (
    name,
    state_id
)
VALUES
    (
        'Bengaluru Urban',
        (
            SELECT id
            FROM states
            WHERE name = 'Karnataka'
        )
    ),
    (
        'Chitradurga',
        (
            SELECT id
            FROM states
            WHERE name = 'Karnataka'
        )
    ),
    (
        'Davangere',
        (
            SELECT id
            FROM states
            WHERE name = 'Karnataka'
        )
    ),
    (
        'Mysuru',
        (
            SELECT id
            FROM states
            WHERE name = 'Karnataka'
        )
    ),
    (
        'Kochi',
        (
            SELECT id
            FROM states
            WHERE name = 'Kerala'
        )
    );

--------------------------------------------------------------------------------------------------------------
-- PINCODES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO pincodes (
    code,
    district_id
)
VALUES
    (
        '560001',
        (
            SELECT id
            FROM districts
            WHERE name = 'Bengaluru Urban'
        )
    ),
    (
        '577501',
        (
            SELECT id
            FROM districts
            WHERE name = 'Chitradurga'
        )
    ),
    (
        '577002',
        (
            SELECT id
            FROM districts
            WHERE name = 'Davangere'
        )
    ),
    (
        '570001',
        (
            SELECT id
            FROM districts
            WHERE name = 'Mysuru'
        )
    );

--------------------------------------------------------------------------------------------------------------
-- DISCIPLINES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO disciplines (
    code,
    display_name
)
VALUES
    ('RIFLE', 'Rifle'),
    ('PISTOL', 'Pistol'),
    ('SHOTGUN', 'Shotgun'),
    ('PARA', 'Para');

--------------------------------------------------------------------------------------------------------------
-- ACADEMIES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO academies (
    name,
    district_id,
    address
)
VALUES
    (
        'Eagle Shooting Academy',
        (
            SELECT id
            FROM districts
            WHERE name = 'Bengaluru Urban'
        ),
        'MG Road, Bengaluru'
    ),
    (
        'Precision Shooting Club',
        (
            SELECT id
            FROM districts
            WHERE name = 'Mysuru'
        ),
        'Vijayanagar, Mysuru'
    ),
    (
        'Champion Rifle Academy',
        (
            SELECT id
            FROM districts
            WHERE name = 'Davangere'
        ),
        'BTM Layout, Davangere'
    );

--------------------------------------------------------------------------------------------------------------
-- USERS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO users (
    first_name,
    last_name,
    email,
    password_hash,
    contact_number,
    role_id
)
VALUES

-- Super Admin
(
    'Bhaskar',
    'VB',
    'bhaskarvb28@gmail.com',
    '$2a$10$NXPrHEt7doY0e2Y1FgyoX.YAwP81O3UGN8Kn1/8MPwPbWxu.sSl1K',
    '+919620122786',
    (
        SELECT id
        FROM roles
        WHERE code = 'super_admin'
    )
),

-- State Admin
(
    'State',
    'Admin',
    'stateadmin@test.com',
    -- State123
    '$2a$10$ue48./bQIuwWnkXRAox3iO1oZqcWXGjaNETeaW/IH11A586fE5ZL.',
    '+919999999991',
    (
        SELECT id
        FROM roles
        WHERE code = 'state_admin'
    )
),

-- District Admin
(
    'District',
    'Admin',
    'districtadmin@test.com',
    '$2a$10$Rq0ww8mjOYNfMnLRjMXOa.dPEZk935clbJqqKFTqEaC1clG7ZPTtm',
    '+919999999992',
    (
        SELECT id
        FROM roles
        WHERE code = 'district_admin'
    )
),

-- District Coach
(
    'District',
    'Coach',
    'districtcoach@test.com',
    '$2a$10$Rq0ww8mjOYNfMnLRjMXOa.dPEZk935clbJqqKFTqEaC1clG7ZPTtm',
    '+919999999993',
    (
        SELECT id
        FROM roles
        WHERE code = 'district_coach'
    )
),

-- Academy Admin
(
    'Academy',
    'Admin',
    'academyadmin@test.com',
    '$2a$10$hk0yiizP5Rr8T4UeD1YMr.C013oyXyvn9HK9qQ5c3aHtRA1rmlPQm',
    '+919999999994',
    (
        SELECT id
        FROM roles
        WHERE code = 'academy_admin'
    )
),

-- Academy Coach
(
    'Academy',
    'Coach',
    'academycoach@test.com',
    '$2a$10$hk0yiizP5Rr8T4UeD1YMr.C013oyXyvn9HK9qQ5c3aHtRA1rmlPQm',
    '+919999999995',
    (
        SELECT id
        FROM roles
        WHERE code = 'academy_coach'
    )
),

-- Player
(
    'Academy',
    'Pllayer',
    'player@test.com',
    '$2a$10$KJ1l1YPpGcB/ClRwPUlkZOsuCSxc.j0FDk5xaKr4QD.1yuoTebhju',
    '+919999999996',
    (
        SELECT id
        FROM roles
        WHERE code = 'player'
    )
);

--------------------------------------------------------------------------------------------------------------
-- STATE ADMINS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO state_admins (
    user_id,
    state_id,
    dpdp_consent,
    profile_completed
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'stateadmin@test.com'
    ),
    (
        SELECT id
        FROM states
        WHERE name = 'Karnataka'
    ),
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT ADMINS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO district_admins (
    user_id,
    district_id,
    dpdp_consent,
    profile_completed
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'districtadmin@test.com'
    ),
    (
        SELECT id
        FROM districts
        WHERE name = 'Bengaluru Urban'
    ),
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACHES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO district_coaches (
    user_id,
    district_id,
    coach_code,
    coaching_certificate_proof,
    dpdp_consent,
    profile_completed
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'districtcoach@test.com'
    ),
    (
        SELECT id
        FROM districts
        WHERE name = 'Bengaluru Urban'
    ),
    'DC001',
    'district_coach_certificate.pdf',
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- DISTRICT COACH DISCIPLINES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO district_coach_disciplines (
    coach_user_id,
    discipline_id
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'districtcoach@test.com'
    ),
    (
        SELECT id
        FROM disciplines
        WHERE code = 'RIFLE'
    )
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY ADMINS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO academy_admins (
    user_id,
    academy_id,
    gstin,
    registration_proof,
    dpdp_consent,
    profile_completed
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'academyadmin@test.com'
    ),
    (
        SELECT id
        FROM academies
        WHERE name = 'Eagle Shooting Academy'
    ),
    '29ABCDE1234F1Z5',
    'academy_registration.pdf',
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACHES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO academy_coaches (
    user_id,
    academy_id,
    coach_code,
    coaching_certificate_proof,
    dpdp_consent,
    profile_completed
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'academycoach@test.com'
    ),
    (
        SELECT id
        FROM academies
        WHERE name = 'Eagle Shooting Academy'
    ),
    'AC001',
    'academy_coach_certificate.pdf',
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- ACADEMY COACH DISCIPLINES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO academy_coach_disciplines (
    coach_user_id,
    discipline_id
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'academycoach@test.com'
    ),
    (
        SELECT id
        FROM disciplines
        WHERE code = 'PISTOL'
    )
);

--------------------------------------------------------------------------------------------------------------
-- PLAYERS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO players (
    user_id,
    academy_id,
    current_coach_user_id,
    status,
    dpdp_consent,
    profile_completed,
    registered_by
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    (
        SELECT id
        FROM academies
        WHERE name = 'Eagle Shooting Academy'
    ),
    (
        SELECT id
        FROM users
        WHERE email = 'academycoach@test.com'
    ),
    'active',
    TRUE,
    TRUE,
    (
        SELECT id
        FROM users
        WHERE email = 'academyadmin@test.com'
    )
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER DISCIPLINES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_disciplines (
    player_user_id,
    discipline_id,
    is_primary
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    (
        SELECT id
        FROM disciplines
        WHERE code = 'PISTOL'
    ),
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER SPORTS PROFILE SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_sports_profile (
    player_user_id,
    unit_of_representation,
    dominant_hand,
    height_cm,
    weight_kg,
    shoe_size,
    tracksuit_size
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    'Karnataka',
    'right',
    175.50,
    68.20,
    '9',
    'M'
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER PERSONAL INFO SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_personal_info (
    player_user_id,
    date_of_birth,
    gender,
    nationality,
    city,
    residential_address,
    pincode_id,
    education,
    institution_name,
    occupation,
    temporary_sport_id
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    '2005-06-15',
    'male',
    'Indian',
    'Bengaluru',
    'Indiranagar, Bengaluru',
    (
        SELECT id
        FROM pincodes
        WHERE code = '560001'
    ),
    'B.Tech',
    'VTU',
    'Student',
    'TMP1001'
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER GUARDIANS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_guardians (
    player_user_id,
    full_name,
    relationship,
    contact_number,
    parental_consent,
    is_primary
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    'Rajesh Singh',
    'Father',
    '+919888888888',
    TRUE,
    TRUE
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER PASSPORTS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_passports (
    player_user_id,
    passport_number,
    passport_issue_date,
    passport_expiry_date,
    passport_issuing_authority,
    passport_place_of_issue
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    'P1234567',
    '2022-01-01',
    '2032-01-01',
    'Government of India',
    'Bengaluru'
);

--------------------------------------------------------------------------------------------------------------
-- PLAYER COACH ASSIGNMENTS SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO player_coach_assignments (
    player_user_id,
    coach_user_id,
    assigned_from,
    is_active
)
VALUES (
    (
        SELECT id
        FROM users
        WHERE email = 'player@test.com'
    ),
    (
        SELECT id
        FROM users
        WHERE email = 'academycoach@test.com'
    ),
    CURRENT_DATE,
    TRUE
);