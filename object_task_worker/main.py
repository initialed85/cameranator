import os

from .consumer import run


def main():
    host = os.getenv("AMQP_HOST") or "localhost"
    port = int(os.getenv("AMQP_PORT") or "5672")
    userid = os.getenv("AMQP_USERID") or "guest"
    password = os.getenv("AMQP_PASSWORD") or "guest"

    run(
        host=host,
        port=port,
        userid=userid,
        password=password,
    )
