## Running the example

The example works with any postgres uri passed as the first argument, we can make this easy by running postgres in a container:

```sh
docker run --name example-db -e POSTGRES_PASSWORD=password -p 8432:5432 -d postgres
go run . postgresql://postgres:password@127.0.0.1:8432
```

You can then jump into a postgres shell to play around with it:
```sh
docker exec -it example-db psql -U postgres
```
