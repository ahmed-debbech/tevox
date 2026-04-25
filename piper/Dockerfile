FROM python:3.12-slim

ARG VOICE_NAME
ARG VOICE_QUALITY
ARG VOICE_LANGUAGE
ARG VOICE_LANGUAGE_FAMILY

ENV VOICE_NAME=${VOICE_NAME}
ENV VOICE_QUALITY=${VOICE_QUALITY}
ENV VOICE_LANGUAGE=${VOICE_LANGUAGE}
ENV VOICE_LANGUAGE_FAMILY=${VOICE_LANGUAGE_FAMILY}

WORKDIR /v

RUN apt update && apt install -y wget 

RUN pip install piper-tts

RUN wget "https://huggingface.co/rhasspy/piper-voices/resolve/main/${VOICE_LANGUAGE_FAMILY}/${VOICE_LANGUAGE}/${VOICE_NAME}/${VOICE_QUALITY}/${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx?download=true" -O ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx
RUN wget "https://huggingface.co/rhasspy/piper-voices/resolve/main/${VOICE_LANGUAGE_FAMILY}/${VOICE_LANGUAGE}/${VOICE_NAME}/${VOICE_QUALITY}/${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx.json?download=true" -O ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx.json

RUN mkdir -p /output

RUN mkdir -p /input

CMD ["sh", "-c", "piper -m ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx --input-file /input/out.txt --output_file /output/hello.wav"]
