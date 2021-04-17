CREATE TABLE IF NOT EXISTS tribe_kms_keyset
(
    id          varchar(255)  NOT NULL,
    created_at  timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    encrypted_keyset          TEXT NOT NULL,
    description varchar(255)  NULL,
    CONSTRAINT pk_tribe_kms_keyset PRIMARY KEY (id)
);
