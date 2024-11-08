CREATE TABLE IF NOT EXISTS users (
    guid uuid PRIMARY KEY NOT NULL,
    login varchar(20) NOT NULL,
    password varchar(60) NOT NULL
);

CREATE UNIQUE INDEX users_login_uniq_idx ON users (login);

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.guid IS 'GUID';
COMMENT ON COLUMN users.login IS 'Логин';
COMMENT ON COLUMN users.password IS 'Пароль';