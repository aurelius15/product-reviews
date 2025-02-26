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
    first_name text     not null,
    last_name  text     not null,
    comment    text     not null,
    rating     smallint not null,
    product_id bigint   not null
        constraint fk_reviews_product
            references products
            on delete cascade,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
