create table products
(
    id          bigserial
        primary key,
    name        text    not null,
    description text    not null,
    price       numeric not null,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone
);

create table reviews
(
    id         bigserial
        primary key,
    first_name text,
    last_name  text,
    comment    text,
    rating     smallint,
    product_id bigint
        constraint fk_products_reviews
            references products
            on delete cascade,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);