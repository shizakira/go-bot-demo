create table if not exists users
(
    id bigint generated always as identity
);

create table if not exists telegram_users
(
    id          bigint generated always as identity,
    user_id     bigint,
    telegram_id bigint,
    chat_id     bigint,
    username    varchar(255),
    foreign key (user_id) references users (id)
);

create table if not exists tasks
(
    id          bigint generated always as identity,
    user_id     int not null,
    title       varchar(255) not null ,
    description text,
    done        bool not null default false,
    deadline    timestamp not null,
    created_at  timestamp not null default now(),
    closed_at   timestamp default null,
    foreign key (user_id) references users (id)
);