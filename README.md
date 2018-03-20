1. In the global mysql configuration file `/etc/my.cnf` or the `.my.cnf` file in your home directory, add the following lines:
```sh
[mysqld]
server-id=master
binlog-format=ROW
log-bin=~/.mysql-binlogs/binlog
```
2. Create a user with replication permissions
```mysql
CREATE USER 'spreader' IDENTIFIED BY 'spreader'; # User a proper password
GRANT REPLICATION SLAVE ON *.* TO 'spreader';
```
3.- Put your AWS credentials in a file located in `~/.aws/credentials`. The contents of the file should be:
```sh
[default]
aws_access_key_id = <Your key>
aws_secret_access_key = <Your secret>
```
