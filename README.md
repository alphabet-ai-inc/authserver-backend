# authserver
Authentication Server for any application with the configuration of applications and databases access

Instructions for go to production

a) Front-end: it is necessary to set up .env and .env.local

b) Back-end:
- see changes in middleware
- compile: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapps ./cmd/api

c) Database:
pg_dump --no-owner -h localhost -p 5432 -U<POSTGRES_USER> autserver > autserver.sql

It is necessary only to copy the front end, go.apps, go.sum and autserver.sql


sudo lsof -i :8080