alter table evaluations
    add constraint user_tea_unique unique (user_id, tea_id);

alter table evaluations
    drop constraint evaluations_rating_check;

alter table evaluations
    add constraint evaluations_rating_check  check ( rating > 0 and rating <= 10);

alter table users
add column name varchar(255),
add column lastname varchar(255);