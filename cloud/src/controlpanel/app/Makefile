
all: clean

clean: 
	find . -name 'controlpanel' -type f -exec rm -f {} \; 
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
	go fmt controlpanel

install:
	go install controlpanel

build:
	go build controlpanel

start:
	./controlpanel 

docker-build: 
	docker build -t controlpanel .
	docker images

docker-run:
	docker run --name controlpanel -td -p 3000:3000 controlpanel
	docker ps

docker-network:
	docker network inspect host

docker-shell:
	docker exec -it controlpanel bash 

docker-clean:
	docker stop controlpanel
	docker rm controlpanel
	docker rmi controlpanel

docker-ip:
	docker-machine ip

heroku-tag:
	docker build -t controlpanel .
	docker tag controlpanel registry.heroku.com/pnguyen-controlpanel/web

heroku-push:
	heroku login 
	heroku container:login
	docker push registry.heroku.com/pnguyen-controlpanel/web

test-ping:
	curl localhost:3000/ping

test-get-inventory:
	curl localhost:3000/controlpanel

test-update-inventory:
	curl -X PUT \
  	http://localhost:3000/controlpanel \
  	-H 'Content-Type: application/json' \
  	-d '{ \
  		"CountGumballs": 1000 }' 

test-place-order:
	curl -X POST \
  	http://localhost:3000/order \
  	-H 'Content-Type: application/json'

test-order-status:
	curl -X GET \
  	http://localhost:3000/order \
  	-H 'Content-Type: application/json'

test-process-order:
	curl -X POST \
  	http://localhost:3000/orders \
  	-H 'Content-Type: application/json'
