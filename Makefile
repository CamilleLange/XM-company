# build only source code
build:
	go build -a -mod=vendor -ldflags '-linkmode external -w -s -extldflags "-static"' -o ./xm-company

# PRODUCTION operations #

build_image:
	docker build -t xm-company:v1.0.0 .

# STAGING operations #

build_image_staging:
	docker build -t xm-company:staging .

deploy_staging: build_image_staging
	docker compose -f deployment/docker-compose-staging.yaml up -d --build

# DEV operations #

build_image_dev:
	docker build -t xm-company:dev .

deploy_dev: build_image_dev
	docker compose -f deployment/docker-compose-dev.yaml up -d --build