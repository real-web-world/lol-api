CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS trigger AS
$$
BEGIN
    NEW.utime = clock_timestamp();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

create table public.config
(
    id    bigserial
        constraint config_pk
            primary key,
    k     varchar(30)                           not null,
    v     text,
    ctime timestamptz default current_timestamp not null,
    utime timestamptz default current_timestamp not null
);

comment on table public.config is '配置';
comment on column public.config.id is 'id';
comment on column public.config.k is '键';
comment on column public.config.v is '值';
comment on column public.config.ctime is '创建时间';
comment on column public.config.utime is '更新时间';

CREATE TRIGGER update_timestamp_config
    BEFORE UPDATE
    ON config
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();