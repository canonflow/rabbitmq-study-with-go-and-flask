FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \
    dos2unix \
    && rm -rf /var/lib/apt/lists/*

COPY flask/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY flask .

RUN dos2unix .env

COPY scripts/flask-start.sh ./
RUN chmod +x ./flask-start.sh 

EXPOSE 3000

CMD ["./flask-start.sh"]