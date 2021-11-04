CREATE TABLE events
(
    id          bigserial not null primary key,
    user_id     int       not null,
    title       varchar   not null,
    description text      not null,
    time        timestamp not null,
    timezone    varchar   not null,
    duration    int8      not null,
    notes       json,
    CONSTRAINT events_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE
);