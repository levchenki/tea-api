alter table teas
    rename column unit_price to weight_price;


alter table teas
    drop column unit_id;

drop table if exists units;

drop type if exists weight_unit;

