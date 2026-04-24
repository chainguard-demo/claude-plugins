FROM cgr.dev/chainguard/python:latest-dev@sha256:2c0fbbac86b72ebb4bfee15b64d8cd5fd6b49dfe7bb279b5c9f193198a84c1c9

WORKDIR /app

COPY requirements.txt .

USER root
RUN pip install --no-cache-dir -r requirements.txt
USER 65532

COPY . .

EXPOSE 8000

ENTRYPOINT ["python", "app.py"]
