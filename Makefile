mysql:
	docker run --name mysql_urlshortener -p 2252:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0.31
createdb:
	docker exec -i mysql_urlshortener mysql -uroot -psecret -e "create database url_short"
dropdb:
	docker exec -i mysql_urlshortener mysql -uroot -psecret -e "drop database url_short"
importdb:
	docker exec -i mysql_urlshortener mysql -uroot -psecret url_short < ./db/sql/query.sql
test:
	go test -v -cover ./...

.PHONY: mysql createdb dropdb importdb test