FROM python:3.11-slim

RUN groupadd -r advisory && useradd -r -g advisory advisory
RUN pip install flask==3.0.0 gunicorn==21.2.0 requests==2.31.0 redis==5.0.0
WORKDIR /app
COPY app /app
COPY cmd.sh /

EXPOSE 9090
USER advisory

CMD ["/cmd.sh"]
