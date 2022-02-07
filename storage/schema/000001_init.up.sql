CREATE TABLE users
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    password_hash bytea,
    gender        varchar(10)  not null,
    age           int          not null,
    start_date    date         default now(),
    description   varchar(255) not null,
    user_hash     bytea
);
 