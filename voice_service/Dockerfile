FROM tevox_vase:latest

WORKDIR /v

RUN apt update && apt install -y wget python3 python3-pip


RUN pip install --break-system-packages piper-tts

ARG VOICE_NAME
ARG VOICE_QUALITY
ARG VOICE_LANGUAGE
ARG VOICE_LANGUAGE_FAMILY

ENV VOICE_NAME=${VOICE_NAME}
ENV VOICE_QUALITY=${VOICE_QUALITY}
ENV VOICE_LANGUAGE=${VOICE_LANGUAGE}
ENV VOICE_LANGUAGE_FAMILY=${VOICE_LANGUAGE_FAMILY}

RUN wget "https://huggingface.co/rhasspy/piper-voices/resolve/main/${VOICE_LANGUAGE_FAMILY}/${VOICE_LANGUAGE}/${VOICE_NAME}/${VOICE_QUALITY}/${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx?download=true" -O ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx
RUN wget "https://huggingface.co/rhasspy/piper-voices/resolve/main/${VOICE_LANGUAGE_FAMILY}/${VOICE_LANGUAGE}/${VOICE_NAME}/${VOICE_QUALITY}/${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx.json?download=true" -O ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx.json

RUN mkdir -p /output

RUN mkdir -p /input

# the app logic

WORKDIR /app

COPY ./voice_service .

CMD ["/go1.24.4/go/bin/go", "run", "."]

#CMD ["sh", "-c", "piper -m ${VOICE_LANGUAGE}-${VOICE_NAME}-${VOICE_QUALITY}.onnx --input-file /input/out.txt --output_file /output/hello.wav"]
