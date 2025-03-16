begin;

create table if not exists categories
(
    id          uuid default gen_random_uuid() primary key,
    name        varchar not null unique,
    description varchar null
);

create table if not exists teas
(
    id          uuid                    default gen_random_uuid() primary key,
    name        varchar        not null unique,
    price       numeric(10, 2) not null,
    description varchar        null,
    created_at  timestamp      not null default current_timestamp,
    updated_at  timestamp      not null default current_timestamp,
    is_deleted  boolean        not null default false,
    category_id uuid           not null references categories (id)
);

create table if not exists tags
(
    id    uuid default gen_random_uuid() primary key,
    name  varchar not null unique,
    color varchar(7)
);

create table if not exists teas_tags
(
    id     uuid default gen_random_uuid() primary key,
    tag_id uuid references tags (id),
    tea_id uuid references teas (id)
);

create table if not exists users
(
    id          uuid                    default gen_random_uuid() primary key,
    telegram_id int unique     not null,
    username    varchar unique null,
    phone       varchar unique,
    created_at  timestamp      not null default current_timestamp,
    updated_at  timestamp      not null default current_timestamp
);

create table if not exists evaluations
(
    id         uuid                                        default gen_random_uuid() primary key,
    rating     numeric(5, 2) check ( rating > 0 ) not null,
    note       varchar                            null,
    created_at timestamp                          not null default current_timestamp,
    updated_at timestamp                          not null default current_timestamp,
    tea_id     uuid references teas (id),
    user_id    uuid references users (id)
);

commit;