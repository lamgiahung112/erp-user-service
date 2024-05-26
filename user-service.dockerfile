FROM alpine:latest

RUN mkdir /app

COPY userServiceApp /app
COPY ./utils/ /app/

CMD ["/app/userServiceApp"]