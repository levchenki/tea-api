create extension if not exists pg_trgm;

create index if not exists idx_teas_name_tgrm on teas using gist (lower(name) gist_trgm_ops);
alter database tea_api_db set pg_trgm.similarity_threshold = 0.1;
