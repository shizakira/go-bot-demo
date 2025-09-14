create table if not exists users
(
    id serial primary key
);

create table if not exists telegram_users
(
    id       serial primary key,
    user_id  bigint,
    chat_id  bigint,
    username varchar(255),
    foreign key (user_id) references users (id)
);

create table if not exists tasks
(
    id          serial primary key,
    user_id     int,
    title       varchar(255),
    description text,
    deadline    timestamp,
    foreign key (user_id) references users (id)
);