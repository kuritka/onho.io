#!/bin/bash
xterm -title "coordinator service 1" -e  "go run main.go coordinator --name=c0 -c=amqp://guest:guest@localhost:5672" &
xterm -title "coordinator service 2" -e  "go run main.go coordinator --name=c1 -c=amqp://guest:guest@localhost:5672" &
xterm -title "sensor-mock service 1" -e  "go run main.go sensor-mock --name=AAA -c=amqp://guest:guest@localhost:5672 -f=1" &
xterm -title "sensor-mock service 2" -e  "go run main.go sensor-mock --name=BBB -c=amqp://guest:guest@localhost:5672 -f=2" &
xterm -title "sensor-mock service 3" -e  "go run main.go sensor-mock --name=CCC -c=amqp://guest:guest@localhost:5672 -f=4" &
xterm -title "database service" -e "go run main.go dbmgr --name=databaseManager -c=amqp://guest:guest@localhost:5672"