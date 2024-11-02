CREATE TABLE IF NOT EXISTS users (
    guid uuid PRIMARY KEY NOT NULL,
    email varchar(50) NOT NULL,
    password varchar(100) NOT NULL
);

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.guid IS 'GUID';
COMMENT ON COLUMN users.email IS 'Email';
COMMENT ON COLUMN users.password IS 'Пароль';