alter table evaluations
    drop constraint user_tea_unique;

alter table evaluations
    drop constraint evaluations_rating_check;

alter table evaluations
    add constraint evaluations_rating_check check ( rating > 0);


alter table users
    drop column name,
    drop column lastname;