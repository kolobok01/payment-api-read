FROM openjdk:11-alpine
ENTRYPOINT ["/usr/bin/payment-api-read.sh"]

COPY payment-api-read.sh /usr/bin/payment-api-read.sh
COPY target/payment-api-read.jar /usr/share/payment-api-read/payment-api-read.jar
