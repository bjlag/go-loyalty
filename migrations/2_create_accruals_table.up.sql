CREATE TABLE IF NOT EXISTS accruals (
    order_number varchar(50) NOT NULL PRIMARY KEY,
    user_guid uuid NOT NULL REFERENCES users (guid),
    status smallint NOT NULL DEFAULT 0,
    accrual integer NOT NULL DEFAULT 0,
    uploaded_at timestamp with time zone
);

CREATE INDEX accruals_user_guid_fk_idx ON accruals (user_guid);

COMMENT ON TABLE accruals IS 'Начисления по заказам';
COMMENT ON COLUMN accruals.order_number IS 'Номер заказа';
COMMENT ON COLUMN accruals.user_guid IS 'GUID пользователя';
COMMENT ON COLUMN accruals.status IS 'Статус обработки начисления по заказу';
COMMENT ON COLUMN accruals.accrual IS 'Размер начисления по заказу';
COMMENT ON COLUMN accruals.uploaded_at IS 'Дата и время загрузки заказа';