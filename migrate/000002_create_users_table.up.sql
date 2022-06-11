CREATE TABLE IF NOT EXISTS users
(
    user_id     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,

    mail        VARCHAR(254) UNIQUE                             NOT NULL,

    secret_type auth_type                                       NOT NULL,

    secret      VARCHAR(163)                                    NOT NULL
);