#!/bin/bash
xterm -title "coordinator 1" -e  "go run main.go coordinator --name=c0 -c=amqp://guest:guest@localhost:5672" &
xterm -title "coordinator 2" -e  "go run main.go coordinator --name=c1 -c=amqp://guest:guest@localhost:5672" &
xterm -title "sensor-mock 1" -e  "go run main.go sensor-mock --name=AAA -c=amqp://guest:guest@localhost:5672 -f=1" &xterm -title "App 1" -e
xterm -title "sensor-mock 2" -e  "go run main.go sensor-mock --name=BBB -c=amqp://guest:guest@localhost:5672 -f=2" &
xterm -title "sensor-mock 3" -e  "go run main.go sensor-mock --name=CCC -c=amqp://guest:guest@localhost:5672 -f=4"