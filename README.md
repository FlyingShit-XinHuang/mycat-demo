# Overview

This is a demo to demonstrate how to query Mycat in a Go program with an ORM package called 'gorm'.

## Mycat environment

Docker image could be built according to [this](https://github.com/FlyingShit-XinHuang/docker-mycat).

### Mysql startup

Start mysql1 and mysql2 containers with following command:

```
docker run -e MYSQL_ROOT_PASSWORD=123456 -d --rm -v `pwd`/my.cnf:/etc/my.cnf --name mysql1 mysql:5.7

docker run -e MYSQL_ROOT_PASSWORD=123456 -d --rm -v `pwd`/my.cnf:/etc/my.cnf --name mysql1 mysql:5.7
```

The my.cnf is mounted so that 'lower_case_table_names=1' option which is needed by Mycat could be added. 
It may look like:

```
[mysqld]

sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES 
lower_case_table_names=1
```

Then create a demo database on each Mysql node:

```
docker exec -ti mysql1 mysql -p123456 -e "create database demo"

docker exec -ti mysql2 mysql -p123456 -e "create database demo"
```

### Mycat configuration

Configure Mycat users in server.xml so that clients could connect to Mycat with them:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!-- - - Licensed under the Apache License, Version 2.0 (the "License"); 
	- you may not use this file except in compliance with the License. - You 
	may obtain a copy of the License at - - http://www.apache.org/licenses/LICENSE-2.0 
	- - Unless required by applicable law or agreed to in writing, software - 
	distributed under the License is distributed on an "AS IS" BASIS, - WITHOUT 
	WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. - See the 
	License for the specific language governing permissions and - limitations 
	under the License. -->
<!DOCTYPE mycat:server SYSTEM "server.dtd">
<mycat:server xmlns:mycat="http://io.mycat/">
    <!-- 
        Create a 'demo' user with password 'demo'. 
        Logical database 'demo' could be queried with this user.
    -->
	<user name="demo">
		<property name="password">demo</property>
		<property name="schemas">demo</property>
	</user>

</mycat:server>
```

Configure a logical database 'demo' and two logical table 'whispir_spaces' and 'whispir_templates' in schema.xml:

```xml
<?xml version="1.0"?>
<!DOCTYPE mycat:schema SYSTEM "schema.dtd">
<mycat:schema xmlns:mycat="http://io.mycat/">
    <!-- Two logical tables in a logical data base 'demo' -->
	<schema name="demo" checkSQLschema="false" sqlMaxLimit="100" >
	    <!-- Both tables are sharded to two data node -->
		<table name='whispir_spaces' primaryKey='ID' autoIncrement="true" rule="mod-long" dataNode="demo1,demo2" />
		<table name='whispir_templates' primaryKey='ID' autoIncrement="true" rule="spaces-mod-long" dataNode="demo1,demo2" />
	</schema>

    <!-- A data node is a database at a Mysql server -->
	<dataNode name="demo1" dataHost="mysql1" database="demo" />
	<dataNode name="demo2" dataHost="mysql2" database="demo" />

    <!-- Configurations of connection to Mysql server -->
	<dataHost name="mysql1" maxCon="1000" minCon="2" balance="0"
	   writeType="0" dbType="mysql" dbDriver="native">
	   <heartbeat>select 1</heartbeat>
	   <writeHost host="mysql1M1" url="172.17.0.2:3306" user="root" password="123456" />
	</dataHost>

	<dataHost name="mysql2" maxCon="1000" minCon="2" balance="0"
	   writeType="0" dbType="mysql" dbDriver="native">
	   <heartbeat>select 1</heartbeat>
	   <writeHost host="mysql2M1" url="172.17.0.3:3306" user="root" password="123456" />
	</dataHost>
</mycat:schema>
```

The IP of mysql1 and mysql2 containers could be obtained with commands:

```
$ docker inspect mysql1|grep IPAddr
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.2",
                    "IPAddress": "172.17.0.2",
```

Configure sharding rules in rule.xml:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mycat:rule SYSTEM "rule.dtd">
<mycat:rule xmlns:mycat="http://io.mycat/">
    <!-- Configure sharding algorithm and the colume to which the algorithm will be applied -->
	<tableRule name="mod-long">
		<rule>
			<columns>id</columns>
			<algorithm>mod-long</algorithm>
		</rule>
	</tableRule>
	<tableRule name="spaces-mod-long">
		<rule>
			<columns>space_id</columns>
			<algorithm>mod-long</algorithm>
		</rule>
	</tableRule>
	
	<!-- Configure sharding function and options -->
	<function name="mod-long" class="io.mycat.route.function.PartitionByMod">
		<!-- how many data nodes -->
		<property name="count">2</property>
	</function>
</mycat:rule>
```

The 'whispir_spaces' table records are sharded using 'mod-long' rule with primary key 'id' and the 'whispir_templates' 
are sharded according to the 'space_id' colume which is a foreign key referencing primary key of 'whispir_spaces'. So all
records of the same space in the two table could reside in a same data node according to above configurations. 

### Create table

Start Mycat container with configurations mounted so that they could be modified more easily:

```
docker run -ti --rm -v `pwd`/conf/:/usr/local/mycat/conf/ -p 8066:8066 -p 9066:9066 mycat:1.6
```

Connect to Mycat and create tables:

```mysql
CREATE TABLE `whispir_spaces` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `whispir_templates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `content` blob NOT NULL,
  `space_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `space_id` (`space_id`),
  CONSTRAINT `whispir_templates_ibfk_1` FOREIGN KEY (`space_id`) REFERENCES `whispir_spaces` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

Global auto-increment id should be configured according to [this article](./auto-increment.md) before any records could 
be inserted.

## Run 

The Go program will do some basic queries with gorm package. Run the demo with main.go.