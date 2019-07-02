# onho.io


## installation


[information about kubernetes deployment](https://stackoverflow.com/questions/25540711/docker-postgres-pgadmin-local-connection)

### RabbitMQ
```bash
sudo docker run -d --hostname my-rabbit --name some-rabbit \
-e RABBITMQ_DEFAULT_USER=guest \
-e RABBITMQ_DEFAULT_PASS=guest \
-p 15672:15672 -p 5671:5671 \
-p 5672:5672 -p 15671:15671 \
-p 25672:25672 rabbitmq:3-management
```

than connection string to mq is `amqp://guest:guest@localhost:5672` and default Admin console is `http://localhost:15672 guest@guest`

[documentation](https://godoc.org/github.com/streadway/amqp) 

### PLSQL

**If you didn't connect yet ensure that `$HOME/docker/volumes/postgres` is empty!**
Database is persited on `$HOME/docker/volumes/postgres` so if you run docker second time
do not remove this directory!

```bash
sudo docker run --rm --name pg-docker \
 -e POSTGRES_PASSWORD=password \
 -p 5432:5432 \
 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data \
 -d postgres:latest
```


If you require custom db and user:
```bash
sudo docker run --rm --name pg-docker \
 -e POSTGRES_PASSWORD=password \
 -e POSTGRES_USER=test \
 -e POSTGRES_DB=test_db \
 -p 5432:5432 \
 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data \
 -d postgres:latest
```


for troubleshooting run `telnet 127.0.0.1 22` from your machine
or `sudo docker exec -it <hash> bash` and check environment variables 
by `env` command and login into database by `psql -h 127.0.0.1 -p 5432 -U test -d postgres --password`
`\l` shows list of available databases


When you create scheme you could export scheme this way... 
```bash
sudo docker exec -it 5f31020172b8 pg_dump --schema-only --no-owner distributed -U test > create_the_tables.sql
``` 

### PGAdmin-4

```bash
sudo docker run -p 8085:80 \
--link pg-docker:pg-docker \
-e "PGADMIN_DEFAULT_EMAIL=admin@local.com" \
-e "PGADMIN_DEFAULT_PASSWORD=password" \
-d dpage/pgadmin4
```

than run `localhost:8085` from your browser.
When create new server use `pg-docker` as Host and `5432` as port 
`postgres` as database and user. For password use `password`



To run application from console you can use mocks:

``xx.sh``
```bash
#!/bin/bash
go run main.go coordinator --name=c0 -c=amqp://guest:guest@localhost:5672 &
go run main.go coordinator --name=c1 -c=amqp://guest:guest@localhost:5672 &
go run main.go sensor-mock --name=AAA -c=amqp://guest:guest@localhost:5672 -f=1 &
go run main.go sensor-mock --name=BBB -c=amqp://guest:guest@localhost:5672 -f=2 &
go run main.go sensor-mock --name=CCC -c=amqp://guest:guest@localhost:5672 -f=4 
```