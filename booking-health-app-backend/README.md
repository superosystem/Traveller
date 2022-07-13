# Backend for Booking Health Application
This repository for handle backend on Booking Health Application with Spring.

## Development
Setup database via docker:
```yml
version: '2.6'
services:
  db:
    container_name: 'healthapp'
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_DATABASE: 'healthapp'
      MYSQL_USER: 'healthapp'
      MYSQL_PASSWORD: 'DZg2JVG2K7NQ7kW0XEUx'
      MYSQL_ROOT_PASSWORD: 'DZg2JVG2K7NQ7kW0XEUx'
    ports:
      - '3306:3306'
    expose:
      - '3306'
    volumes:
      - "${PWD}/my_db:/var/lib/mysql"

volumes:
  my_db:
```
Running docker compose
```bash
docker-compose up
```