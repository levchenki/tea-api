create table users_favourite_teas
(
    id         uuid                                default gen_random_uuid() primary key,
    user_id    uuid references users (id) not null,
    tea_id     uuid references teas (id)  not null,
    created_at timestamp                  not null default current_timestamp,
    constraint favourite_user_tea_unique unique (user_id, tea_id)
);