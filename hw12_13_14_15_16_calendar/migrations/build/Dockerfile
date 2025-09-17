FROM gomicro/goose

ADD ./migrations/*.sql /migrations/
ADD ./migrations/entrypoint.sh /migrations/

RUN chmod +x /migrations/entrypoint.sh



ENTRYPOINT ["/migrations/entrypoint.sh"]