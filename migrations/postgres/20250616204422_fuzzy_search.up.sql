create extension if not exists pg_trgm;

create index if not exists idx_teas_name_tgrm on teas using gist (lower(name) gist_trgm_ops);

DO
$$
    BEGIN
        EXECUTE format('ALTER DATABASE %I SET pg_trgm.similarity_threshold = 0.1;', current_database());
    END
$$;

alter table teas
    rename column is_deleted to is_hidden;