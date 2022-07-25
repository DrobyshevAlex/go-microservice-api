FROM rabbitmq:3-management

ENV DOCKER_RABBITMQ_USER dev
ENV DOCKER_RABBITMQ_USER dev
ENV DOCKER_RABBITMQ_USER /var/lib/rabbitmq/mnesia/rabbitmq

ADD init.sh /init.sh
RUN chmod +x /init.sh

RUN apt-get update
RUN apt-get install -y curl 
EXPOSE 4369 5671 5672 25672 15671 15672

CMD ["/init.sh"]
