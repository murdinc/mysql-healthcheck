# mysql-healthcheck

> Mysql Health Check Daemon

## Intro

**mysql-healthcheck** is a simple program designed to run as a daemon to respond to ELB health checks for mysql servers. It checks DB connection and runs a test query to determine if an instance is healthy.

## Configuration

The configuration file is loaded from: `/etc/mysql-healthcheck.cnf` and the options are very simple to configure. Example:

```cnf
[Mysql]
port = 3306
username= user
password = pw
database = db

[HealthCheck]
port = 9999
```
