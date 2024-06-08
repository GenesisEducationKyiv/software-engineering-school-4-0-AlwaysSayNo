CREATE SEQUENCE users_id_seq;

CREATE TABLE users
(
    id         BIGINT DEFAULT nextval('users_id_seq') NOT NULL PRIMARY KEY,
    email      VARCHAR(255)                           NOT NULL UNIQUE
);