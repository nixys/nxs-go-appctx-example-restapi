# nxs-go-appctx-example-restapi

This projects implements the following elements:
- [nxs-go-appctx](https://github.com/nixys/nxs-go-appctx) library helps bring together all the components for create a complete REST API server
- Project structure you may use for it (described in [nxs-go-appctx](https://github.com/nixys/nxs-go-appctx) documentation)
- [GORM](https://github.com/go-gorm/gorm) library to work with databases (MySQL in this case)
- [Gin Web Framework](github.com/gin-gonic/gin) to create a REST API
- [golang-migrate](https://github.com/golang-migrate/migrate) to write the database migrations

## Description

The application represents an API server with implemented CRUD operations for manage user accounts stored in MySQL.

## Quickstart

### Requrements

To run the test application make sure you have installed following tools:
- [Docker Compose](https://docs.docker.com/compose/install/)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

### Prepare environment

- Go to application directory
- Start MySQL Database via Docker Compose:
  ```bash
  docker compose -f .env/docker-compose.yml up -d
  ```
- Run database migration:
  ```bash
  migrate -path migrations/ -database mysql://user:somepass@tcp\(127.0.0.1:3306\)/db?multiStatements=true up
  ```

### Run test app

- Run API server with following command executed within the application directory:
```bash
go run main.go
```
- To test application make a following requests to API:
  - Create new user account:
    ```bash
    curl -X POST -H "X-Auth-Key: some_auth_token" -d '{"username": "userA"}'  "http://127.0.0.1:8080/v1/user"
    ```
    Response:
    ```json
    {
      "user": {
        "id": 1,
        "username": "userA",
        "password": "ettih20vId4JOMq"
      }
    }
    ```
  - List existing accounts:
    ```bash
    curl -X GET -H "X-Auth-Key: some_auth_token" "http://127.0.0.1:8080/v1/user"
    ```
    Response:
    ```json
    {
      "users": [
        {
          "id": 1,
          "username": "userA",
          "password": "ettih20vId4JOMq"
        }
      ]
    }
    ```
  - Update account with ID: 1:
    ```bash
    curl -X PATCH -H "X-Auth-Key: some_auth_token" -d '{"password": "fixed password"}'  "http://127.0.0.1:8080/v1/user/1"
    ```
    Response:
    ```json
    {
      "user": {
        "id": 1,
        "username": "userA",
        "password": "fixed password"
      }
    }
    ```
  - Get account with ID: 1:
    ```bash
    curl -X GET -H "X-Auth-Key: some_auth_token"  "http://127.0.0.1:8080/v1/user/1"
    ```
    Response:
    ```json
    {
      "user": {
        "id": 1,
        "username": "userA",
        "password": "fixed password"
      }
    }
    ```
  - Delete account with ID: 1:
    ```bash
    curl -X DELETE -H "X-Auth-Key: some_auth_token"  "http://127.0.0.1:8080/v1/user/1"
    ```

## Feedback

For support and feedback please contact me:
- telegram: [@borisershov](https://t.me/borisershov)
- e-mail: b.ershov@nixys.ru

## License

nxs-go-appctx-example-restapi is released under the [MIT License](LICENSE).
