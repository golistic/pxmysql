/*
 * Copyright (c) 2022, Geert JM Vanderkelen
 */

USE pxmysql_tests;

DROP TABLE IF EXISTS `data_types_datetime`;

CREATE TABLE data_types_datetime
(
    id           TINYINT AUTO_INCREMENT,
    dt_date      DATE NOT NULL,
    dt_time      TIME(6) NOT NULL,
    dt_datetime  DATETIME(6) NOT NULL,
    dt_timestamp TIMESTAMP NOT NULL,
    dt_year      YEAR NOT NULL,
    PRIMARY KEY (id)
);

SET @@time_zone = '+00:00';
INSERT INTO data_types_datetime
VALUES (1, '2005-03-01', '08:00:01.123456', '2005-03-01 07:00:01', FROM_UNIXTIME(1109660401),
        2005),
       (2, '9999-12-31', '838:59:59.0', '9999-12-31 23:59:59.999999', FROM_UNIXTIME(2147483647),
        1901),
       (3, '1000-01-01', '-838:59:59.0', '1000-01-01 00:00:00', FROM_UNIXTIME(1),
        1901);
