USE pxmysql_tests;

DROP TABLE IF EXISTS `data_types_string`;

CREATE TABLE `data_types_string`
(
    id          TINYINT AUTO_INCREMENT,
    s_char      CHAR(255) NOT NULL,
    s_varchar   VARCHAR(400) NOT NULL,
    s_binary    BINARY(20) NOT NULL,
    s_varbinary VARBINARY(20) NOT NULL,
    s_longtext  LONGTEXT NOT NULL,
    s_tinyblob  TINYBLOB NOT NULL,
    s_enum      ENUM ('Go', 'Python', 'JavaScript') NOT NULL,
    s_set       SET ('Go', 'Python', 'JavaScript') NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO data_types_string
VALUES (1,
        CONCAT('CHAR', REPEAT('a', 251)),
        CONCAT('VARCHAR', REPEAT('b', 393)),
        X'0708090a0b0c0d0e0f10',
        X'08090a0b0c0d0e0f10',
        CONCAT('LONGTEXT', REPEAT('l', @@mysqlx_max_allowed_packet - 10)),
        'I am a tiny blob',
        'Go',
        'Python,Go');
