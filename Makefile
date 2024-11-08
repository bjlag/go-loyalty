up:
	docker run -d --rm \
		--name loyalty_pg \
		-p 5432:5432 \
		-v ${PWD}/data/psql:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=secret \
		-e POSTGRES_DB=master \
		postgres:16.4-alpine3.20

down:
	docker stop loyalty_pg

exec:
	docker exec -it loyalty_pg psql -U postgres