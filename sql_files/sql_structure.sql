create table postcodes_source
(
    postcode        varchar(10)  null,
    postcode_spaces varchar(10)  null,
    in_use          varchar(5)   null,
    country         varchar(50)  null,
    county          varchar(150) null,
    constraint uk_postcode
        unique (postcode)
);

