# axdb

This project was created to train myself in Go.

## Starting the server

```sh
make run
```

## Usage examples

### REST API
1. PUT data
    ```sh
    curl -H 'Content-Type: application/json' -X PUT -d '{"value":"example value"}' http://localhost:8080/items/example
    ```

2. GET data
    ```sh
    curl -H 'Content-Type: application/json' -X GET http://localhost:8080/items/example
    ```

3. GET data index
    ```sh
    curl -H 'Content-Type: application/json' -X GET http://localhost:8080/items
    ```