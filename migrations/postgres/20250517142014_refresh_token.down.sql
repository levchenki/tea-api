alter table users
    drop column refresh_token_id;

alter table teas
    rename column serve_price to price;

alter table teas
    drop column weight_price;