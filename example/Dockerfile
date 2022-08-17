ARG IMAGE_REGISTRY
ARG IMAGE_CACHE_TAG
ARG IMAGE_TAG

FROM ${IMAGE_REGISTRY}zaaksam/golang:${IMAGE_CACHE_TAG} AS build-api-service
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o gins_example

FROM ${IMAGE_REGISTRY}amibaio/alpine:${IMAGE_TAG}
LABEL maintainer="zaaksam"
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY helper/bin/ /tmp/
COPY --from=build-api-service /app/gins_example /app/
RUN chmod +x /usr/local/bin/entrypoint.sh \
    && chmod +x /tmp/install.sh \
    && sh /tmp/install.sh
EXPOSE 8080
ENTRYPOINT ["entrypoint.sh"]