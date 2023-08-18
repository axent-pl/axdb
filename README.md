# axdb

This project was created to train myself in Go.

## Usage examples

1. REST PUT data
    ```sh
    curl -H 'Content-Type: application/json' -X PUT -d '{"name":"John", "comment":"is fine"}' http://localhost:6600/items/example
    ```

2. REST GET data
    ```sh
    curl -H 'Content-Type: application/json' -X GET http://localhost:6600/items/example
    ```

3. REST GET data index
    ```sh
    curl -H 'Content-Type: application/json' -X GET http://localhost:6600/items
    ```