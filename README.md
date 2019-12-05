# go-ms-poc

## Authentication

- OAuth2 (JWT) with support for keycloak
- See: [AUTH.md](AUTH.md)

## Checklist

- [x] Health EP
  - Checks have to be implemented on our own (e.g. state of DB)
- [x] Version EP
- [x] OAuth + Role Check (using keycloak)
  - PoC impl. working
  - Based on github.com/auth0-community/go-auth0 and gopkg.in/square/go-jose.v2
- [x] PostgreSQL
  - PoC impl. working
  - Based on database/sql and github.com/lib/pq
  - Connection Pooling/ Handling has to be implemented on our own
- [ ] RabbitMQ
- [ ] Graph DB (neo4j)
- [x] Logging (structured)
- [x] Config (ENV + CLI)
- [x] Docker MS
- [ ] REST
- [ ] GraphQL
- [x] Graceful Shutdown