# Auth

## Local Setup

### Get the keycloak docker image

```bash
docker pull quay.io/keycloak/keycloak
```

### Start keycloak

```bash
docker run -d \
    -e KEYCLOAK_USER=admin \
    -e KEYCLOAK_PASSWORD=admin \
    -p 8180:8080 \
    -v `pwd`/testdata/goms-realm.json:/tmp/goms-realm.json \
    -e KEYCLOAK_IMPORT=/tmp/goms-realm.json \
    --name goms-keycloak \
    quay.io/keycloak/keycloak

# open keycloak in browser
xdg-open http://localhost:8180

# start it again
docker stop goms-test
```

### Stop keycloak & cleanup

```bash
# stop it
docker stop goms-test

# cleanup
docker rm goms-test
```
