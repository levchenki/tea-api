alter table teas
    rename column weight_price to unit_price;

create type weight_unit as enum ('G', 'KG');

create table if not exists units
(
    id          uuid default gen_random_uuid() primary key,
    is_apiece   bool default false not null,
    weight_unit weight_unit        not null,
    value       int                not null
);

alter table teas
    add column unit_id uuid references units (id);

