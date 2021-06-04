ARG GOLANG=1.15

# Build stage
FROM golang:${GOLANG}-alpine AS BUILD-STAGE
LABEL stage=build
RUN apk add --no-cache curl git openssh-client build-base
WORKDIR /app
COPY . /app
RUN make all

# Final stage WEB
FROM alpine
WORKDIR /app
RUN apk --no-cache add curl
RUN apk --no-cache upgrade
COPY --from=BUILD-STAGE /app/output/ /app/
COPY --from=BUILD-STAGE /app/entrypoint.sh /app/entrypoint.sh
