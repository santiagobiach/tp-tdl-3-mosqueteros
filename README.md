Por ahora:

- para correr el servidor. en la carpeta server:

    go build -o server main.go
    ./server 12000

- para correr un cliente. en la carpeta principal:

    go build -o client client.go
    ./client localhost:12000
