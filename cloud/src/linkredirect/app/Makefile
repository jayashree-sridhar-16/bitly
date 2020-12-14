
all: clean

clean: 
	find . -name 'linkredirect' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/codegangsta/negroni
	go get -v github.com/gorilla/mux
	go get -v github.com/unrolled/render
	go get -v github.com/satori/go.uuid
	go get -v github.com/go-sql-driver/mysql

format:
	go fmt linkredirect

install:
	go install linkredirect

build:
	go build linkredirect

start:
	./linkredirect 

docker-build: 
	docker build -t linkredirect .
	docker images

docker-run:
	docker run --name linkredirect -td -p 3000:3000 linkredirect
	docker ps

docker-network:
	docker network inspect host

docker-shell:
	docker exec -it linkredirect bash 

docker-clean:
	docker stop linkredirect
	docker rm linkredirect
	docker rmi linkredirect

docker-ip:
	docker-machine ip

heroku-tag:
	docker build -t linkredirect .
	docker tag linkredirect registry.heroku.com/pnguyen-linkredirect/web

heroku-push:
	heroku login 
	heroku container:login
	docker push registry.heroku.com/pnguyen-linkredirect/web

test-ping:
	curl localhost:3000/ping

test-get-inventory:
	curl localhost:3000/linkredirect

test-update-inventory:
	curl -X PUT \
  	http://localhost:3000/linkredirect \
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
