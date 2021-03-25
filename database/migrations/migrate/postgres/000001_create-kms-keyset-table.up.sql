CREATE TABLE IF NOT EXISTS tribe_kms_keyset
(
    id          varchar(255)  NOT NULL,
    created_at  timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    name        varchar(255)  NOT NULL,
    next_id     varchar(255)  NULL,
    encrypted_keyset          TEXT NOT NULL,
    description varchar(255)  NULL,
    CONSTRAINT tribe_kms_keyset_id_pk PRIMARY KEY (id)
);
