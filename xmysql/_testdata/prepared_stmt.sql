/*
 * Copyright (c) 2022, Geert JM Vanderkelen
 */

USE pxmysql_tests;

DROP TABLE IF EXISTS `numeric_not_null`, `numeric_null`;
DROP TABLE IF EXISTS `temporal_not_null`, `temporal_null`;
DROP TABLE IF EXISTS `strings_not_null`, `strings_null`;

CREATE TABLE numeric_not_null
(
    id                 TINYINT AUTO_INCREMENT,
    bit_               BIT(6)             NOT NULL,
    bool_              BOOL               NOT NULL,
    tinyint_           TINYINT            NOT NULL,
    tinyint_unsigned   TINYINT UNSIGNED   NOT NULL,
    smallint_          SMALLINT           NOT NULL,
    smallint_unsigned  SMALLINT UNSIGNED  NOT NULL,
    mediumint_         MEDIUMINT          NOT NULL,
    mediumint_unsigned MEDIUMINT UNSIGNED NOT NULL,
    int_               INT                NOT NULL,
    int_unsigned       INT UNSIGNED       NOT NULL,
    bigint_            BIGINT             NOT NULL,
    bigint_unsigned    BIGINT UNSIGNED    NOT NULL,
    decimal_           DECIMAL(65, 30)    NOT NULL,
    float_             FLOAT              NOT NULL,
    float_unsigned     FLOAT UNSIGNED     NOT NULL,
    double_            DOUBLE             NOT NULL,
    double_unsigned    DOUBLE UNSIGNED    NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE numeric_null
(
    id                 TINYINT AUTO_INCREMENT,
    bit_               BIT(6)             NULL,
    bool_              BOOL               NULL,
    tinyint_           TINYINT            NULL,
    tinyint_unsigned   TINYINT UNSIGNED   NULL,
    smallint_          SMALLINT           NULL,
    smallint_unsigned  SMALLINT UNSIGNED  NULL,
    mediumint_         MEDIUMINT          NULL,
    mediumint_unsigned MEDIUMINT UNSIGNED NULL,
    int_               INT                NULL,
    int_unsigned       INT UNSIGNED       NULL,
    bigint_            BIGINT             NULL,
    bigint_unsigned    BIGINT UNSIGNED    NULL,
    decimal_           DECIMAL(65, 30)    NULL,
    float_             FLOAT              NULL,
    float_unsigned     FLOAT UNSIGNED     NULL,
    double_            DOUBLE             NULL,
    double_unsigned    DOUBLE UNSIGNED    NULL,
    PRIMARY KEY (id)
);

CREATE TABLE temporal_not_null
(
    id         TINYINT AUTO_INCREMENT,
    datetime_  DATETIME(6)  NOT NULL,
    date_      DATE         NOT NULL,
    timestamp_ TIMESTAMP(6) NOT NULL,
    year_      YEAR         NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE temporal_null
(
    id         TINYINT AUTO_INCREMENT,
    datetime_  DATETIME(6)  NULL,
    date_      DATE         NULL,
    timestamp_ TIMESTAMP(6) NULL,
    year_      YEAR         NULL,
    PRIMARY KEY (id)
);

CREATE TABLE strings_not_null
(
    id          TINYINT AUTO_INCREMENT,
    char_       CHAR(255)                                NOT NULL,
    binary_     BINARY(255)                              NOT NULL,
    varchar_    VARCHAR(600)                             NOT NULL,
    varbinary_  VARBINARY(410)                           NOT NULL,
    tinyblob_   TINYBLOB                                 NOT NULL,
    tinytext_   TINYTEXT                                 NOT NULL,
    blob_       BLOB                                     NOT NULL,
    text_       TEXT                                     NOT NULL,
    mediumblob_ MEDIUMBLOB                               NOT NULL,
    mediumtext_ MEDIUMTEXT                               NOT NULL,
    longblob_   LONGBLOB                                 NOT NULL,
    longtext_   LONGTEXT                                 NOT NULL,
    enum_       ENUM ('Earth', 'Moon', 'Mars', 'Europa') NOT NULL,
    set_        SET ('Earth', 'Moon', 'Mars')            NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE strings_null
(
    id          TINYINT AUTO_INCREMENT,
    char_       CHAR(255)      NULL,
    binary_     BINARY(255)    NULL,
    varchar_    VARCHAR(600)   NULL,
    varbinary_  VARBINARY(410) NULL,
    tinyblob_   TINYBLOB       NULL,
    tinytext_   TINYTEXT       NULL,
    blob_       BLOB           NULL,
    text_       TEXT           NULL,
    mediumblob_ MEDIUMBLOB     NULL,
    mediumtext_ MEDIUMTEXT     NULL,
    longblob_   LONGBLOB       NULL,
    longtext_   LONGTEXT       NULL,
    enum_       ENUM ('Earth', 'Moon', 'Mars'),
    set_        SET ('Earth', 'Moon', 'Mars'),
    PRIMARY KEY (id)
);