CREATE SEQUENCE currencies_id_seq;

CREATE TABLE currencies
(
    id     BIGINT DEFAULT nextval('currencies_id_seq') NOT NULL PRIMARY KEY,
    number FLOAT                                       NOT NULL,
    date   VARCHAR(64)                                 NOT NULL
);