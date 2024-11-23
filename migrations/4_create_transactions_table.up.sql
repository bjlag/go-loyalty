CREATE TABLE IF NOT EXISTS transactions (
    guid uuid NOT NULL PRIMARY KEY,
    account_guid uuid NOT NULL REFERENCES accounts (guid),
    order_number varchar(50) NOT NULL,
    type smallint NOT NULL,
    sum double precision NOT NULL,
    processed_at timestamp with time zone NOT NULL
);

COMMENT ON TABLE transactions IS 'Транзакции по счету';
COMMENT ON COLUMN transactions.guid IS 'GUID транзакции';
COMMENT ON COLUMN transactions.account_guid IS 'GUID счета';
COMMENT ON COLUMN transactions.order_number IS 'Номер заказа, по которому зарегистрирована транзакция';
COMMENT ON COLUMN transactions.type IS 'Тип транзакции: 0 - начислили, 1 - списали';
COMMENT ON COLUMN transactions.sum IS 'Сумма транзакции';
COMMENT ON COLUMN transactions.processed_at IS 'Дата и время, когда прошла транзакция';
