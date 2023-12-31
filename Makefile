.PHONY: db
db:
	go build -o db ./cmd/database/main.go
	./db

.PHONY: server
server:
	go build -o server ./cmd/server/main.go
	./server
 
.PHONY: publisher
publisher:
	go build -o publisher ./cmd/publisher/main.go
	./publisher

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: clean
clean:
	rm db server publisher
