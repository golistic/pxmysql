USE pxmysql_tests;

DROP TABLE IF EXISTS not_collection_28380dew22;
DROP TABLE IF EXISTS collection_wic28skwixkd;
DROP TABLE IF EXISTS collection_weux73293jsnsj;

CREATE TABLE not_collection_28380dew22
(
    id INT
);

CREATE TABLE collection_wic28skwixkd
(
    doc          json null,
    _id          varbinary(32) as (json_unquote(json_extract(`doc`, _utf8mb4'$._id'))) stored
        primary key,
    _json_schema json as (_utf8mb4'{"type":"object"}'),
    constraint $val_strict_wic28skwixkd
        check (json_schema_valid(`_json_schema`, `doc`))
);

CREATE TABLE collection_weux73293jsnsj
(
    doc          json null,
    _id          varbinary(32) as (json_unquote(json_extract(`doc`, _utf8mb4'$._id'))) stored
        primary key,
    _json_schema json as (_utf8mb4'{"type":"object"}'),
    constraint $val_strict_weux73293jsnsj
        check (json_schema_valid(`_json_schema`, `doc`))
);