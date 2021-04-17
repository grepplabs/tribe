CREATE TABLE IF NOT EXISTS tribe_oidc_jwks
(
    id               varchar(255)  NOT NULL,
    created_at       timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    current_jwks_id  varchar(255)  NOT NULL,
    next_jwks_id     varchar(255)  NOT NULL,
    previous_jwks_id varchar(255)  NULL,
    rotation_mode    integer       NOT NULL DEFAULT 0,
    rotation_period  integer       NOT NULL DEFAULT 0,
    last_rotated     timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    description      varchar(255)  NULL,
    version          integer       NOT NULL DEFAULT 0,
    CONSTRAINT pk_tribe_oidc_jwks PRIMARY KEY (id),
    CONSTRAINT fk_tribe_oidc_jwks_current FOREIGN KEY (current_jwks_id) REFERENCES tribe_jwks (id),
    CONSTRAINT fk_tribe_oidc_jwks_next FOREIGN KEY (next_jwks_id) REFERENCES tribe_jwks (id),
    CONSTRAINT fk_tribe_oidc_jwks_previous FOREIGN KEY (previous_jwks_id) REFERENCES tribe_jwks (id),
    CONSTRAINT tribe_oidc_jwks_diff_jwks_ids CHECK (current_jwks_id != next_jwks_id AND current_jwks_id != previous_jwks_id AND next_jwks_id != previous_jwks_id)
);

CREATE INDEX IF NOT EXISTS tribe_oidc_jwks_current ON tribe_oidc_jwks (current_jwks_id);

CREATE INDEX IF NOT EXISTS tribe_oidc_jwks_next ON tribe_oidc_jwks (next_jwks_id);

CREATE INDEX IF NOT EXISTS tribe_oidc_jwks_previous ON tribe_oidc_jwks (previous_jwks_id);

CREATE INDEX IF NOT EXISTS tribe_oidc_jwks_rotation ON tribe_oidc_jwks (last_rotated, rotation_mode);
