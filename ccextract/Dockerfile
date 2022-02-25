FROM debian:latest 

RUN apt update && apt install -y cargo git libglew-dev libglfw3-dev cmake gcc libcurl4-gnutls-dev tesseract-ocr libtesseract-dev libleptonica-dev clang libclang-dev
RUN git clone https://github.com/CCExtractor/ccextractor.git
WORKDIR ccextractor/linux
RUN ./build -without-rust

CMD ["./ccextractor"]