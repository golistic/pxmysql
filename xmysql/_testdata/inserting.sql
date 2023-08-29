USE pxmysql_tests;

DROP TABLE IF EXISTS `inserts01`;

CREATE TABLE inserts01
(
    id TINYINT NOT NULL AUTO_INCREMENT,
    c1 VARCHAR(20),
    PRIMARY KEY (id)
);
