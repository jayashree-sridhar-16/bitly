
all: clean

clean: 
	find . -name 'trendserver' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/codegangsta/negroni
	go get -v github.com/gorilla/mux
	go get -v github.com/unrolled/render
	go get -v github.com/satori/go.uuid
	go get -v github.com/go-sql-driver/mysql
	go get -v github.com/streadway/amqp

format:
	go fmt trendserver

install:
	go install trendserver

build:
	go build trendserver

start:
	./trendserver 

docker-build: 
	docker build -t trendserver .
	docker images

docker-run:
	docker run --name trendserver -td -p 3000:3000 trendserver
	docker ps

docker-network:
	docker network inspect host

docker-shell:
	docker exec -it trendserver bash 

docker-clean:
	docker stop trendserver
	docker rm trendserver
	docker rmi trendserver

docker-ip:
	docker-machine ip

rabbitmq-run:
	docker run --restart always --name rabbitmq \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

heroku-tag:
	docker build -t trendserver .
	docker tag trendserver registry.heroku.com/pnguyen-trendserver/web

heroku-push:
	heroku login 
	heroku container:login
	docker push registry.heroku.com/pnguyen-trendserver/web

test-ping:
	curl localhost:3000/ping