CREATE TABLE IF NOT EXISTS tribe_realm
(
    realm_id    varchar(64)  NOT NULL,
    description varchar(255) NULL,
    CONSTRAINT tribe_realm_id_pk PRIMARY KEY (realm_id),
    CONSTRAINT tribe_realm_id_pattern CHECK (realm_id ~* '^[a-zA-Z0-9_-]*$')
);
