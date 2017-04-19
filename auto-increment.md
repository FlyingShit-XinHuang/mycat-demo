# How to generate auto increment id with Mycat

This article will demonstrate how to generate global auto increment id in a sharded table.

## Mycat configurations

Configure 'sequnceHandlerType' in server.xml:

```xml
<system>
    <property name="sequnceHandlerType">1</property>
</system>
```

Configure table with 'autoIncrement' option in schema.xml:

```xml
<schema name="demo" checkSQLschema="false" sqlMaxLimit="100" >
	<table name='whispir_users' primaryKey='ID' autoIncrement="true" rule="mod-long" dataNode="demo1,demo2" ></table>
</schema>
```

Auto increment id is implemented with the help of a MYCAT_SEQUENCE table and some functions which are stored in one of 
the Mycat data node. 
Add mapping of table to data node in sequence_db_conf.properties:

```
WHISPIR_USERS=demo1
```

__Note__ the table name in sequence_db_conf.properties must be uppercase.

## Create sequence table and functions

Connect to the Mysql database related with data node 'demo1' and create table and functions as following:

```mysql
DROP TABLE IF EXISTS MYCAT_SEQUENCE;
-- name sequence 名称
-- current_value 当前 value
-- increment 增长步长! 可理解为 mycat 在数据库中一次读取多少个 sequence. 当这些用完后, 下次再从数据库中读取.

CREATE TABLE MYCAT_SEQUENCE (name VARCHAR(50) NOT NULL,current_value INT NOT NULL,increment INT NOT NULL DEFAULT 100, PRIMARY KEY(name)) ENGINE=InnoDB;

-- 获取当前 sequence 的值 (返回当前值,增量)
DROP FUNCTION IF EXISTS mycat_seq_currval;
DELIMITER ;;
CREATE FUNCTION mycat_seq_currval(seq_name VARCHAR(50)) RETURNS varchar(64) CHARSET 'utf8'
DETERMINISTIC
BEGIN
DECLARE retval VARCHAR(64);
SET retval="-999999999,null";
SELECT concat(CAST(current_value AS CHAR),",",CAST(increment AS CHAR)) INTO retval FROM MYCAT_SEQUENCE WHERE name = seq_name;
RETURN retval;
END
;;
DELIMITER ;

-- 设置 sequence 值
DROP FUNCTION IF EXISTS mycat_seq_setval;
DELIMITER ;;
CREATE FUNCTION mycat_seq_setval(seq_name VARCHAR(50),value INTEGER) RETURNS varchar(64) CHARSET 'utf8'
DETERMINISTIC
BEGIN
UPDATE MYCAT_SEQUENCE
SET current_value = value
WHERE name = seq_name;
RETURN mycat_seq_currval(seq_name);
END
;;
DELIMITER ;

-- 获取下一个 sequence 值
DROP FUNCTION IF EXISTS mycat_seq_nextval;
DELIMITER ;;
CREATE FUNCTION mycat_seq_nextval(seq_name VARCHAR(50)) RETURNS varchar(64) CHARSET 'utf8'
DETERMINISTIC
BEGIN
UPDATE MYCAT_SEQUENCE
SET current_value = current_value + increment WHERE name = seq_name;
RETURN mycat_seq_currval(seq_name);
END
;;
DELIMITER ;
```

Insert a record about sequence configuration of table whispir_users:

```mysql
INSERT INTO MYCAT_SEQUENCE(name,current_value,increment) VALUES ('whispir_users', 100000, 100)
```

## Create table

Start Mycat and connect. Create table whose primary key must be AUTO_INCREMENT:

```mysql
create table whispir_users (
  id int not null auto_increment,
  name varchar(100) not null,
  primary key(id)
) engine InnoDB default charset='utf8';
```

Insert records and test:

```mysql
mysql> insert into whispir_users (name) values ('foo');
Query OK, 1 row affected (0.15 sec)

mysql> select last_insert_id();
+------------------+
| LAST_INSERT_ID() |
+------------------+
|           100100 |
+------------------+
```