# Perpustakaan RESTful API

<p align="justify">
This is a RESTful API for Library Management powered by Spring Boot, Java that provide the main function you 
could stock your book and you could borrow a book. This project for final exam when I learn 
at SinauKoding Bootcamp.
</p>

## Technologies
- Java JDK 17
- Maven 3.8.6
- Spring Framework 2.7.1
- MySQL 8.0

## Installation

You have to install Java JDK 17 and MySQL 8.0.
- Clone this repository
```bash
git clone https://github.com/gusrylmubarok/perpustakaan-spring-rest.git
```
- Open on your favorite `IDE` or `Text Editor`
- Run and configuration `MySQL Server`
```text
database name: perpustakaandb
database user: <your_username>
database password: <your_password>
```
- Run and build Application
```bash
mvn clean install package && mvn spring-boot:run
```

## API Reference

#### Registration User
```http
  POST /auth/do-register
```
| Parameter     | Type     | Description                     |
|:--------------| :------- |:--------------------------------|
| `profileName` | `string` | **Required**. Your Profile Name |
| `username`    | `string` | **Required**. Your Username     |
| `password`    | `string` | **Required**. Your Password     |

#### Login User
```http
  POST /auth/do-login
```
| Parameter     | Type     | Description                     |
|:--------------| :------- |:--------------------------------|
| `username`    | `string` | **Required**. Your Username     |
| `password`    | `string` | **Required**. Your Password     |

#### Login User
```http
  POST /auth/do-login
```
| Parameter     | Type     | Description                     |
|:--------------| :------- |:--------------------------------|
| `username`    | `string` | **Required**. Your Username     |
| `password`    | `string` | **Required**. Your Password     |

#### Book API
#### Loan API
#### User API


## TODO
- [x] Authentication(Register and Login)
- [ ] CRUD and Search Books
- [ ] CRUD and Search Loan

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://github.com/gusrylmubarok/perpustakaan-spring-rest/blob/main/LICENSE.md)
