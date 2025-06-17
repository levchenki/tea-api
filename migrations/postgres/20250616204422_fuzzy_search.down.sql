drop index if exists idx_teas_name_tgrm;

alter database tea_api_db reset pg_trgm.similarity_threshold;
drop extension if exists pg_trgm;
