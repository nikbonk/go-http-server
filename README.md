```bash
./startPostgres.sh
```

```bash
psql -h localhost -p 5432 -U postgres
```

```bash
export GOOSE_MIGRATION_DIR="$(pwd)/sql/schema"
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
```
