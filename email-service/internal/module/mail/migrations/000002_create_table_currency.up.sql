CREATE SEQUENCE currency_id_seq;

CREATE TABLE currency
(
    id    BIGINT DEFAULT nextval('currency_id_seq') NOT NULL PRIMARY KEY,
    value BIGINT                                    NOT NULL
);