alter table users
    add column refresh_token_id uuid null;

alter table teas
    rename column price to serve_price;

alter table teas
    add column weight_price numeric(10, 2) null;

update teas
set weight_price = serve_price;


alter table teas
    alter column weight_price set not null;