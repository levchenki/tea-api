alter table evaluations
    drop constraint user_tea_unique;

alter table evaluations
    drop constraint evaluations_rating_check;

alter table evaluations
    add constraint evaluations_rating_check check ( rating > 0);


alter table users
    drop column first_name,
    drop column last_name,
    drop column is_admin,
    add column phone varchar(255),
    alter column telegram_id type int using telegram_id::int;
