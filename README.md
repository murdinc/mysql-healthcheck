# mysql-healthcheck
> Mysql Health Check Daemon

[![Build Status](https://travis-ci.org/murdinc/mysql-healthcheck.svg)](https://travis-ci.org/murdinc/mysql-healthcheck)


## Intro
**mysql-healthcheck** is a A simple program designed to run as a daemon to respond to ELB health checks for mysql servers. It checks the query count and slave status to determine if an instance is healthy.


## Installation
```
curl -s http://dl.sudoba.sh/get/mysql-healthcheck | sh
```

## Configuration
The configuration file is loaded from: `/etc/mysql/mysql-healthcheck.cnf` and the options are very simple to configure. Example:

```
[Mysql]
port = 3306
username= user
password = pw
database = db

[HealthCheck]
port = 9999
maxqueries = 150
checkslavestatus # make sure that mysql slave is running

```