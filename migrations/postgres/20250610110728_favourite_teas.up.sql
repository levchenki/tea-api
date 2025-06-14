create table users_favourite_teas
(
    id         uuid               default gen_random_uuid() primary key,
    user_id    uuid references users (id),
    tea_id     uuid references teas (id),
    created_at timestamp not null default current_timestamp
);