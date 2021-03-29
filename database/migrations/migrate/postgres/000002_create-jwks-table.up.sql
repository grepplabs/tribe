CREATE TABLE IF NOT EXISTS tribe_jwks
(
    id               varchar(255)  NOT NULL,
    created_at       timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    kid              varchar(255)  NOT NULL,
    alg              varchar(32)   NOT NULL,
    use              varchar(32)   NOT NULL,
    kms_keyset_uri   varchar(255)  NOT NULL,
    encrypted_jwks   TEXT NOT NULL,
    description      varchar(255)  NULL,
    CONSTRAINT pk_tribe_jwks PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS tribe_jwks_kid_use ON tribe_jwks (kid,use);
