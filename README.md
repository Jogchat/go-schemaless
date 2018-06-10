This is an open-source, MIT-licensed implementation of Uber's Schemaless
(immutable BigTable-style sharded MySQL datastore)

All code is in Golang, no exceptions.

## DATABASE SUPPORT

For learning or other:

	* SQLite (the 'fs' and 'memory' storages are just file and memory
	  SQLite backends)

	* rqlite (Distributed SQLite) - experimental, broken

For potentially serious usage:

	* MySQL

	* Postgres


## ADDING SUPPORT FOR ADDITIONAL DATABASES / STORAGES

I will be more than happy to accept well-tested, high-quality implementations
for other potential storage backends. If you need support for something (but
don't see it here) then please file an issue to open discussion up. PRs
welcome.

## SETTING UP FOR DEVELOPMENT AND RUNNING TESTS

1. Install MySQL, postgres, and rqlite, setup users on MySQL and Postgres.

2 Run both shell scripts inside tools/create_shard_schemas, one at a time,
loading the generated sql file into Postgres and MySQL locally.

(TODO: Future versions will split the tool's functionality in a way
that it can be integrated into an organization's "data fabric", creating
schemas + grants semi-automatically for you, on different sets of shards)

3. Now, you can run tests a bit more easily. For me, this looks like:

~/go-src/src/github.com/rbastic/go-schemaless$ MYSQLUSER=user MYSQLPASS=pass PGUSER=user PGPASS=pass SQLHOST=localhost make test

Having replaced the user and pass with the appropriate usernames and passwords
for MySQL and Postgres, this should pass all tests.

Any test cases should be idempotent - they should not result in errors on
subsequent runs due to hard-coded row keys.

## DISCLAIMER

I do not work for Uber Technologies. Everything has been sourced from their
materials that they've released on the subject matter (which I am extremely
gracious for): 

## VIDEOS

"Taking Storage for a Ride With Uber", https://www.youtube.com/watch?v=Dg76cNaeB4s (30 mins)

"GOTO 2016 • Taking Storage for a Ride", https://www.youtube.com/watch?v=kq4gp90QUcs (1 hour)

## ARTICLES

"Designing Schemaless, Uber Engineering’s Scalable Datastore Using MySQL"

"Part One", https://eng.uber.com/schemaless-part-one/

"Part Two", https://eng.uber.com/schemaless-part-two/

"Part Three", https://eng.uber.com/schemaless-part-three/

"Code Migration in Production: Rewriting the Sharding Layer of Uber’s Schemaless Datastore"
https://eng.uber.com/schemaless-rewrite/

The underlying sharding code is https://github.com/dgryski/go-shardedkv/choosers,
similar versions of which have powered https://github.com/photosrv/photosrv and
also a large sharded MySQL database system. The storage and storagetest code is
also derived from https://github.com/dgryski/go-shardedkv

My sincere thanks to Damian Gryski for open-sourcing the above package.

## OTHER RESOURCES

FriendFeed: https://backchannel.org/blog/friendfeed-schemaless-mysql

Pinterest: https://engineering.pinterest.com/blog/sharding-pinterest-how-we-scaled-our-mysql-fleet

Martin Fowler's slides on Schemaless Data Structures: https://martinfowler.com/articles/schemaless/

## OTHER OPEN-SOURCE IMPLEMENTATIONS

https://github.com/hoteltonight/shameless - Schemaless in Ruby

## TO LOGIN TO MYSQL SERVER

Connect to MySQL node

	ssh root@165.227.25.43

Connect to MySQL server
    
    mysql -u root -p
    
Have fun

## How to set up Ubuntu 16.04 with schemaless-go server

## Download and install go bin
```
sudo apt-get update
sudo apt-get -y upgrade
sudo curl -O https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
sudo mv go /usr/local
sudo vim ~/.profile
add at end of file:
export PATH=$PATH:/usr/local/go/bin
source ~/.profile
```
## 
```
makedir go
export GOPATH=$HOME/go
cd go
makedir code.jogchat.internal
cd code.jogchat.internal
git clone https://github.com/Jogchat/go_schemaless.git
```

## DESIGN
BLOB DESIGN
```
user blob{
        "id": dummy_id, // uuid
        "username": "dummy_username",
        "email": "dummy_email",
        "password": "dummy_password",
        "activate": False // boolean
}
company blob{
        "id": dummy_id // uuid
        "category": "dummy_category",
        "domain": "dummy_domain",
        "name": "dummy_name"
}
school blob{
        "id": dummy_id // uuid
        "category": "dummy_category"
        "domain": "dummy_domain"
        "name": "dummy_name"
}
```
## How to Allow MySQL Remote Access in Ubuntu Server 16.04
In this tutorial we are going to learn how to allow remote access to the MySQL server in Ubuntu Server. For the tutorial I am using Ubuntu Server 16.04, But you can use this on any previous version of Ubuntu Linux.

Enable MySQL Server Remote Connection in Ubuntu
By default MySQL Server on Ubuntu run on the local interface, This means remote access to the MySQL Server is not Allowed. To enable remote connections to the MySQL Server we need to change value of the bind-address in the MySQL Configuration File.

```
First, Open the /etc/mysql/mysql.conf.d/mysqld.cnf file (/etc/mysql/my.cnf in Ubuntu 14.04 and earlier versions).

vim /etc/mysql/mysql.conf.d/mysqld.cnf
```

Under the [mysqld] Locate the Line,

```
bind-address            = 127.0.0.1
```

And change it to,

```
bind-address            = 0.0.0.0
```

How to Allow MySQL Remote Access in Ubuntu Server 16.04
Then, Restart the Ubuntu MysQL Server.

```
systemctl restart mysql.service
```

Now Ubuntu Server will allow remote access to the MySQL Server, But still you need to configure MySQL users to allow access from any host.

For example, when you create a MySQL user, you should allow access from any host.

```
CREATE USER 'username'@'%' IDENTIFIED BY 'password';
```

Or Allow from Specific IP Address,

```
CREATE USER 'username'@'192.168.1.100' IDENTIFIED BY 'password';
```

Troubleshoot Ubuntu MySQL Remote Access
To make sure that, MySQL server listens on all interfaces, run the netstat command as follows.

```
netstat -tulnp | grep mysql
```

The output should show that MySQL Server running on the socket 0 0.0.0.0:3306 instead of 127.0.0.1:3306.

MySQL Server running on the socket 0 0.0.0.0:3306
You can also try to telnet to the MySQL port 3306 from a remote host. For example, if the IP Address of your Ubuntu Server is 192.168.1.10, Then from the remote host execute,

```
telnet 192.168.1.10 3306
```

You can also run the nmap command from a remote computer to check whether MySQL port 3306 is open to the remote host.

```
nmap 192.168.1.10
```

The output should list MySQL port 3306 and the STATe should be open. If the MySQL port 3306 not open, Then there is a firewall which blocks the port 3306.

Summary : MySQL Remote Access Ubuntu Server 16.04.
In this tutorial we learned how to enable Remote Access to MySQL Server in Ubuntu 16.04.

```
mysql> SHOW GRANTS;
+---------------------------------------------------------------------+
| Grants for root@localhost                                           |
+---------------------------------------------------------------------+
| GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION |
| GRANT PROXY ON ''@'' TO 'root'@'localhost' WITH GRANT OPTION        |
+---------------------------------------------------------------------+
2 rows in set (0.00 sec)
```
To grant user remote access:
```
CREATE USER 'root'@'%' IDENTIFIED BY 'password'
```


