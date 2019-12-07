# go-ms-poc

## Authentication

- OAuth2 (JWT) with support for keycloak
- See: [AUTH.md](AUTH.md)

## Checklist

- [x] Health EP
  - Checks have to be implemented on our own (e.g. state of DB)
  - [ ] health check
- [x] Version EP
- [x] OAuth + Role Check (using keycloak)
  - [x] PoC impl. working
  - Based on github.com/auth0-community/go-auth0 and gopkg.in/square/go-jose.v2
- [x] PostgreSQL
  - [x] PoC impl. working
  - Based on database/sql and github.com/lib/pq
  - Connection Pooling/ Handling has to be implemented on our own
- [x] RabbitMQ
  - [x] PoC impl. working
- [ ] Graph DB (neo4j)
  - [ ] PoC impl. working
- [x] Logging (structured)
- [x] Config (ENV + CLI)
  - [ ] PoC impl. working
- [x] Docker MS
  - [ ] PoC impl. working
- [ ] REST
  - [ ] PoC impl. working
  - http://www.gorillatoolkit.org
    - provides:
      - context handling
      - [x] form to struct conversion (works)
      - securecookies
      - sessionhandling
- [ ] GraphQL
  - [ ] PoC impl. working
- [x] Graceful Shutdown
- [x] Metrics
  - [ ] PoC impl. working
