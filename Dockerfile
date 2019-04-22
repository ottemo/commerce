FROM alpine:3.5

RUN apk add --no-cache ca-certificates gawk

RUN mkdir -pv /home/ottemo/commerce
RUN mkdir -pv /home/ottemo/commerce/var/log
RUN mkdir -pv /home/ottemo/commerce/var/session

COPY commerce /home/ottemo/commerce/

# create links for proper logging
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/cron.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/events.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/paypal.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/rest.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/subscription.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/mongo.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/events.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/models.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/usps.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/impex.log
RUN ln -sf /dev/stdout /home/ottemo/commerce/var/log/product
RUN ln -sf /dev/stderr /home/ottemo/commerce/var/log/errors.log

COPY bin/docker-entrypoint.sh /home/ottemo/commerce/

EXPOSE 3000
WORKDIR /home/ottemo/commerce
CMD ./docker-entrypoint.sh
