import os

from .object_detector import run, run_once, setup

if __name__ == "__main__":
    if os.getenv("SETUP") or "0" == "1":
        setup()
        exit(0)

    if os.getenv("ONE_SHOT") or "0" == "1":
        run_once()
        exit(0)

    run()
