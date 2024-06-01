FROM alpine:latest

RUN mkdir /app

COPY userServiceApp /app
COPY ./external/ /app/

CMD ["/app/userServiceApp"]