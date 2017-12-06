FROM golang
MAINTAINER orchestration
RUN mkdir /app
ADD orchestration /app/
ADD pkcs8.key /app/
ADD cacert.pem /app/
ADD nuvem-intera.local.pem /app/
WORKDIR /app
EXPOSE 5000
CMD ["/app/orchestration"]
