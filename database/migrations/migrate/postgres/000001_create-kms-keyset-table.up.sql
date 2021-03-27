CREATE TABLE IF NOT EXISTS tribe_kms_keyset
(
    id          varchar(255)  NOT NULL,
    name        varchar(255)  NOT NULL,
    created_at  timestamp     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    next_id     varchar(255)  NULL,
    encrypted_keyset          TEXT NOT NULL,
    description varchar(255)  NULL,
    CONSTRAINT pk_tribe_kms_keyset PRIMARY KEY (id,name),
    CONSTRAINT fk_tribe_kms_keyset_next_id FOREIGN KEY (next_id,name) REFERENCES tribe_kms_keyset (id,name)
);

CREATE UNIQUE INDEX IF NOT EXISTS tribe_kms_keyset_name ON tribe_kms_keyset (name) WHERE (next_id is null);

CREATE UNIQUE INDEX IF NOT EXISTS tribe_kms_keyset_next_id ON tribe_kms_keyset (next_id, name) WHERE (next_id is not null);
