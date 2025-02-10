#!/bin/bash
docker stop $(docker ps -aq) 2>/dev/null && docker rm $(docker ps -aq) 2>/dev/null

docker rmi -f $(docker images -q) 2>/dev/null

docker build -t my_go_image .

docker run -d --name my_go_app -p 8080:8080 my_go_image
