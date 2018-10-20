FROM consul
RUN apk update && apk upgrade && apk add bash

RUN wget -O /bin/fabio https://github.com/fabiolb/fabio/releases/download/v1.5.9/fabio-1.5.9-go1.10.2-linux_amd64
RUN chmod +x /bin/fabio

COPY consul-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

