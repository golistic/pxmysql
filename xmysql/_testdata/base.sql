CREATE SCHEMA IF NOT EXISTS pxmysql_tests;

/*
 * Copyright (c) 2022, Geert JM Vanderkelen
 */

-- Following users are used for testing the authentication with
-- MySQL Authentication Plugins.
CREATE USER IF NOT EXISTS 'user_native'@'%' IDENTIFIED WITH mysql_native_password
    BY 'pwd_user_native';
GRANT ALL ON pxmysql_tests.* TO 'user_native'@'%';

CREATE USER IF NOT EXISTS 'user_sha256'@'%' IDENTIFIED WITH caching_sha2_password
    BY 'pwd_user_sha256';
GRANT ALL ON pxmysql_tests.* TO 'user_sha256'@'%';

-- Basic users for testing.
CREATE USER IF NOT EXISTS 'pxmysqltest'@'%' IDENTIFIED WITH mysql_native_password BY '';
GRANT ALL ON pxmysql_tests.* TO 'pxmysqltest'@'%';
