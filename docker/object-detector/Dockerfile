FROM --platform=linux/amd64 pytorch/pytorch:2.2.1-cuda11.8-cudnn8-runtime

RUN apt-get update && apt-get install -y ffmpeg libsm6 libxext6

RUN python3 -m pip install --upgrade pip

WORKDIR /srv/

COPY requirements.txt /srv/requirements.txt

RUN python3 -m pip install -r requirements.txt

COPY object_detector /srv/object_detector

RUN bash -c 'SETUP=1 MODEL_PATH=./object_detector/models/yolov7x.pt python3 -m object_detector'
