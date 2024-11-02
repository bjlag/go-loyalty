CREATE TABLE IF NOT EXISTS users (
    guid uuid PRIMARY KEY NOT NULL,
    email varchar(50) NOT NULL,
    password varchar(60) NOT NULL
);

CREATE UNIQUE INDEX users_email_uniq_idx ON users (email);

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.guid IS 'GUID';
COMMENT ON COLUMN users.email IS 'Email';
COMMENT ON COLUMN users.password IS 'Пароль';