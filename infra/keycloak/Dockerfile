
## The Build Image
FROM quay.io/keycloak/keycloak

USER root

## Copy config files for realms
COPY goms-realm.json /tmp/
COPY keycloak-add-user.json /opt/jboss/keycloak/standalone/configuration/

ENV KEYCLOAK_USER=admin
ENV KEYCLOAK_PASSWORD=admin
ENV KEYCLOAK_IMPORT=/tmp/goms-realm.json

EXPOSE 8080

ENTRYPOINT [ "/opt/jboss/tools/docker-entrypoint.sh" ]

# Default startup including script mappers enabled
CMD  ["-b", "0.0.0.0", "-Dkeycloak.profile.feature.upload_scripts=enabled"]