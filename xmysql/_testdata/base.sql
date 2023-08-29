CREATE SCHEMA IF NOT EXISTS pxmysql_tests;
CREATE SCHEMA IF NOT EXISTS pxmysql_tests_a;

-- Following users are used for testing the authentication with
-- MySQL Authentication Plugins.
CREATE USER IF NOT EXISTS 'user_native'@'%' IDENTIFIED WITH mysql_native_password
    BY 'pwd_user_native';
GRANT ALL ON pxmysql_tests.* TO 'user_native'@'%';
GRANT ALL ON pxmysql_tests_a.* TO 'user_native'@'%';

CREATE USER IF NOT EXISTS 'user_sha256'@'%' IDENTIFIED WITH caching_sha2_password
    BY 'pwd_user_sha256';
GRANT ALL ON pxmysql_tests.* TO 'user_sha256'@'%';
GRANT ALL ON pxmysql_tests_a.* TO 'user_sha256'@'%';

CREATE USER IF NOT EXISTS 'pxmysqltest'@'%' IDENTIFIED WITH mysql_native_password BY '';
GRANT ALL ON pxmysql_tests.* TO 'pxmysqltest'@'%';
GRANT ALL ON pxmysql_tests_a.* TO 'pxmysqltest'@'%';

-- Clean up objects that might have been created by tests
DROP SCHEMA IF EXISTS `pxmysql_2839cks829dka`;
