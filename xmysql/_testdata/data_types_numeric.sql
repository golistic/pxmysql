USE pxmysql_tests;

DROP TABLE IF EXISTS `data_types_numeric`;

CREATE TABLE data_types_numeric
(
    id                         TINYINT AUTO_INCREMENT,
    numeric_bit                BIT(6) NULL,
    numeric_bool               BOOL NULL,
    numeric_tinyint            TINYINT NULL,
    numeric_tinyint_unsigned   TINYINT UNSIGNED NULL,
    numeric_smallint           SMALLINT NOT NULL,
    numeric_smallint_unsigned  SMALLINT UNSIGNED NOT NULL,
    numeric_mediumint          MEDIUMINT NOT NULL,
    numeric_mediumint_unsigned MEDIUMINT UNSIGNED NOT NULL,
    numeric_int                INT NOT NULL,
    numeric_int_unsigned       INT UNSIGNED NOT NULL,
    numeric_bigint             BIGINT NOT NULL,
    numeric_bigint_unsigned    BIGINT UNSIGNED NOT NULL,
    numeric_decimal            DECIMAL(65, 30) NOT NULL,
    numeric_decimal2           DECIMAL(65, 1) NOT NULL,
    numeric_decimal3           DECIMAL(18, 9) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO data_types_numeric
VALUES (1, b'100110', false, 127, 0,
        32767, 0, 8388607, 0, 2147483647, 0,
        9223372036854775807, 0,
        3.14,
        9999999999999999999999999999999999999999999999999999999999991234.9,
        123456789.000001),
       (2, b'000110', true, -128, 255,
        -32768, 65535, -8388608, 16777215, -2147483648, 4294967295,
        -9223372036854775808, 18446744073709551615,
        -3.14,
        -9999999999999999999999999999999999999999999999999999999999991234.5,
        -123456789.000001)
;
