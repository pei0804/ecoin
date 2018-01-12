init:
	@docker-compose run --rm dep init
install:
	@docker-compose run --rm dep ensure
update:
	@docker-compose run --rm dep ensure -update && dep ensure
serve-dev:
	@dev_appserver.py ./app --host 0.0.0.0 --admin_host 0.0.0.0 --skip_sdk_update_check
run:
	@docker-compose up app

docker-server: docker-build docker-up
docker-clean: docker-stop docker-rm

docker-build:
	docker-compose build

docker-up:
	docker-compose up

docker-stop:
	docker-compose stop

docker-rm:
	docker-compose rm

