create table if not exists engines
(
    engine_model  text    not null
        constraint engines_pk
            primary key,
    engine_power  integer not null,
    engine_volume integer,
    engine_type   text    not null
);

alter table engines
    owner to postgres;

create unique index if not exists engines_engine_model_uindex
    on engines (engine_model);

create table if not exists transmissions
(
    transmission_model        text    not null
        constraint transmissions_pk
            primary key,
    transmission_type         text    not null,
    transmission_gears_number integer not null
);

alter table transmissions
    owner to postgres;

create unique index if not exists transmissions_transmission_model_uindex
    on transmissions (transmission_model);

create table if not exists brands
(
    brand_name            text not null
        constraint brands_pk
            primary key,
    brand_creator_country text not null
);

alter table brands
    owner to postgres;

create table if not exists wheels
(
    wheel_radius integer not null,
    wheel_color  text    not null,
    wheel_model  text    not null
        constraint wheels_pk
            primary key
);

alter table wheels
    owner to postgres;

create unique index if not exists wheels_wheel_model_uindex
    on wheels (wheel_model);

create table if not exists cars
(
    model        text    not null
        constraint cars_pk
            primary key,
    brand        text    not null
        constraint cars_brands_brand_name_fk
            references brands
            on update cascade on delete cascade,
    engine       text    not null
        constraint cars_engines_engine_model_fk
            references engines
            on update cascade on delete cascade,
    transmission text    not null
        constraint cars_transmissions_transmission_model_fk
            references transmissions
            on update cascade on delete cascade,
    price        integer not null,
    wheel        text    not null
        constraint cars_wheels_wheel_model_fk
            references wheels
            on update cascade on delete cascade
);

alter table cars
    owner to postgres;

create unique index if not exists cars_model_uindex
    on cars (model);

