mysql-healthcheck
=================

MySQL ELB health check service

A simple program designed to run as a daemon to respond to ELB health checks for mysql servers.

It checks the query count and slave status to determine if an instance is healthy. 
