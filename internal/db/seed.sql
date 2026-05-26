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
-- SUPER ADMIN SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO users (
    first_name,
    last_name,
    email,
    password_hash,
    contact_number,
    role_id
)
VALUES (
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
);

--------------------------------------------------------------------------------------------------------------
-- WEAPON CATEGORIES SEED
--------------------------------------------------------------------------------------------------------------

INSERT INTO weapon_categories (
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