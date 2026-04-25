FROM tevox_vase:latest


ARG TESSERACT_VERSION

ENV TESSERACT_VERSION=${TESSERACT_VERSION}

RUN wget https://github.com/tesseract-ocr/tesseract/archive/refs/tags/${TESSERACT_VERSION}.tar.gz

RUN tar -xvf ${TESSERACT_VERSION}.tar.gz

RUN mv tesseract-${TESSERACT_VERSION} tesseract

WORKDIR tesseract

RUN ./autogen.sh

RUN ./configure

RUN make

RUN make install

RUN ldconfig

RUN wget https://raw.githubusercontent.com/tesseract-ocr/tessdata_best/refs/heads/main/eng.traineddata

RUN mkdir -p /text

RUN mkdir -p /pic

# the app logic

WORKDIR /app

COPY ./ocr_service .

CMD ["/go1.24.4/go/bin/go", "run", "."]
#CMD ["sh", "-c", "tesseract --tessdata-dir ./ /pic/image - -l eng > /text/out.txt"]