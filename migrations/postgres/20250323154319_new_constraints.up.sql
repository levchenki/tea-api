alter table evaluations
    add constraint user_tea_unique unique (user_id, tea_id);

alter table evaluations
    drop constraint evaluations_rating_check;

alter table evaluations
    add constraint evaluations_rating_check check ( rating > 0 and rating <= 10);

alter table users
    add column first_name varchar(255),
    add column last_name  varchar(255),
    add column is_admin   bool default false,
    drop column phone,
    alter column telegram_id type bigint using telegram_id::bigint;
