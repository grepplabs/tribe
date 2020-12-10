CREATE TABLE IF NOT EXISTS tribe_user
(
    user_id            varchar(320) NOT NULL,
    created_at         timestamp    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    realm_id           varchar(64)  NOT NULL,
    username           varchar(255) NOT NULL,
    encrypted_password varchar(255) NOT NULL,
    enabled            boolean      NOT NULL default false,
    email              varchar(320) NULL,
    email_verified     boolean      NOT NULL default false,
    CONSTRAINT tribe_user_id_pk PRIMARY KEY (user_id),
    CONSTRAINT tribe_user_realm_id_fk FOREIGN KEY (realm_id) REFERENCES tribe_realm (realm_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS tribe_user_username_realm_id_idx ON tribe_user (username, realm_id);

CREATE INDEX IF NOT EXISTS tribe_user_email_idx ON tribe_user (email);

CREATE INDEX IF NOT EXISTS tribe_user_realm_id_idx ON tribe_user (realm_id);

CREATE UNIQUE INDEX IF NOT EXISTS tribe_user_email_verified_idx ON tribe_user (email, realm_id) WHERE (email_verified is true);