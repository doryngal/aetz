all: build run

build:
	docker build -t webaetz:1 .

run:
	docker run --rm -d \
		--network host \
		--name webaetz \
  		-e DB_HOST=mypostgres \
  		-e DB_PORT=5432 \
  		-e DB_USER=postgres \
  		-e DB_PASSWORD=postgres \
  		-e DB_NAME=binai \
  		webaetz:1

stop:
	docker stop webaetz
 
restart:
	make all
	sudo systemctl stop nginx
	sudo caddy run --config /etc/caddy/Caddyfile


rebuild:
	make stop
	make build
	make run

runParserHTMLgoszakup:
	cd ./parserHTMLgoszakup \
	go run .
runCommonService:
	cd /home/ubuntu/binai/commonService && go run ./cmd
pq:
	docker exec -it mypostgres psql --host=localhost --dbname=binai --username=baha
# uses in crontab
backup:
	sh /home/ubuntu/binai/createBackup.sh
runVectorModel:
	/usr/bin/python -m jupyter nbconvert --to notebook --execute "path.ipynb"
upLinuxEnv:
	docker-compose -f ./docker-compose.postgresLinux.yml up -d
backup/dl:
	scp -r ubuntu@185.22.67.27:/home/ubuntu/backups/ ./
backup/copy:
	docker cp ./backups mypostgres:/backups
	docker cp ./schemas mypostgres:/schemas
backup/restore:
	# make backup/dl
	make backup/copy
	docker exec -it mypostgres psql --host=localhost --username=postgres -f ./schemas/createRole_baha.sql
	docker exec -it mypostgres psql --host=localhost --username=postgres -f ./schemas/createDatabase_binai.sql
	docker exec -it mypostgres psql --host=localhost --dbname=binai --username=baha -f ./backups/binaiMockBD.sql
backup/createFromDocker:
	docker exec -it mypostgres pg_dump -h localhost -U baha -d binai > ./backups/binaiMockBD.sql
