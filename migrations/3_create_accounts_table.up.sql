CREATE TABLE IF NOT EXISTS accounts (
    guid uuid NOT NULL PRIMARY KEY REFERENCES users (guid),
    balance integer NOT NULL DEFAULT 0,
    withdraw_sum integer NOT NULL DEFAULT 0,
    updated_at timestamp with time zone NOT NULL
);

COMMENT ON TABLE accounts IS 'Счета баллов лояльности пользователя';
COMMENT ON COLUMN accounts.guid IS 'GUID счета, совпадает с GUID пользователя';
COMMENT ON COLUMN accounts.balance IS 'Сумма доступных баллов лояльности на счету';
COMMENT ON COLUMN accounts.withdraw_sum IS 'Сумма использованных за весь период регистрации баллов';
COMMENT ON COLUMN accounts.updated_at IS 'Дата и время обновления счета';
