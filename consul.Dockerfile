FROM consul
RUN apk update && apk upgrade && apk add bash

COPY consul-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

