drop index if exists idx_teas_name_tgrm;

DO
$$
    BEGIN
        EXECUTE format('ALTER DATABASE %I reset pg_trgm.similarity_threshold;', current_database());
    END
$$;

drop extension if exists pg_trgm;
