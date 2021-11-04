CREATE TABLE users
(
    id                 bigserial not null primary key,
    login              varchar   not null unique,
    encrypted_password varchar   not null,
    timezone           varchar   not null
);