create table if not exists users
(
    id          serial,
    telegram_id bigint,
    username    varchar(255)
);

create table if not exists tasks
(
    id          serial,
    title       varchar(255),
    description text,
    deadline    timestamp
);