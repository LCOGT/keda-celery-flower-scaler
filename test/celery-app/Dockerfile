FROM python:3.10-slim

RUN pip install --upgrade pip poetry

WORKDIR /src

COPY . .

RUN poetry install --no-root
