# app.py
import os
import threading
import time
import logging
from flask import Flask, request, jsonify
import json
import pika
from dotenv import load_dotenv

# Load .env for local dev
load_dotenv()

# Configuration
RABBITMQ_URL = os.getenv("RABBITMQ_CONNECTION_URL")
if not RABBITMQ_URL:
    raise RuntimeError("RABBITMQ_CONNECTION_URL must be set")

PUBLISH_QUEUE_NAME = "second_queue"
CONSUME_QUEUE_NAME = "first_queue"

# Logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# ----------------------------
# Consumer Function (runs in background thread)
# ----------------------------
def rabbitmq_consumer():
    while True:
        try:
            params = pika.URLParameters(RABBITMQ_URL)
            connection = pika.BlockingConnection(params)
            channel = connection.channel()

            # Ensure queue exists
            channel.queue_declare(queue=CONSUME_QUEUE_NAME, durable=True)
            channel.basic_qos(prefetch_count=1)

            def callback(ch, method, properties, body):
                logger.info(f" [x] Received: {body.decode()}")
                time.sleep(1)  # Simulate work
                ch.basic_ack(delivery_tag=method.delivery_tag)
                logger.info(" [x] Done")

            channel.basic_consume(queue=CONSUME_QUEUE_NAME, on_message_callback=callback)
            logger.info(" [*] Consumer started. Waiting for messages...")
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            logger.error(f"Consumer connection error: {e}. Retrying in 5s...")
            time.sleep(5)
        except Exception as e:
            logger.error(f"Consumer unexpected error: {e}")
            time.sleep(5)
        finally:
            if 'connection' in locals() and connection.is_open:
                connection.close()

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "OK"}), 200

# ----------------------------
# Flask Routes (Publisher)
# ----------------------------
@app.route('/send', methods=['POST'])
def send_message():
    try:
        # Hard-coded JSON data
        message_data = {
            "source": "flask",
            "message": "Hello From Flask!",
            "timestamp": "2025-10-29T15:30:00Z"
        }

        # Convert to JSON string
        message_body = json.dumps(message_data)

        params = pika.URLParameters(RABBITMQ_URL)
        connection = pika.BlockingConnection(params)
        channel = connection.channel()

        # Declare queue (durable=True is okay, but message won't be persistent)
        channel.queue_declare(queue=PUBLISH_QUEUE_NAME, durable=True)

        # Publish NON-PERSISTENT message
        channel.basic_publish(
            exchange='',
            routing_key=PUBLISH_QUEUE_NAME,
            body=message_body,
            properties=pika.BasicProperties(
                content_type='application/json'  # optional, for metadata only
                # delivery_mode is omitted → defaults to 1 (non-persistent)
            )
        )
        connection.close()

        return jsonify({"status": "sent", "message": message_data}), 200

    except Exception as e:
        return jsonify({"error": str(e)}), 500

# ----------------------------
# Start consumer in a daemon thread on startup
# ----------------------------
def start_consumer_thread():
    consumer_thread = threading.Thread(target=rabbitmq_consumer, daemon=True)
    consumer_thread.start()
    logger.info("Started RabbitMQ consumer thread.")

# Start consumer when module is run directly
if __name__ == '__main__':
    start_consumer_thread()
    app.run(host='0.0.0.0', port=3000, debug=False)