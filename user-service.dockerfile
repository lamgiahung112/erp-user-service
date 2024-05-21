FROM alpine:latest

RUN mkdir /app

COPY userServiceApp /app

CMD ["/app/userServiceApp"]