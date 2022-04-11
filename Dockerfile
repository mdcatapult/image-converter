FROM golang:1.17 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY model ./model
COPY utils ./utils
COPY *.go ./

RUN go mod tidy
RUN go build -o /mdc-minerva-image-converter


FROM fiji/fiji:fiji-openjdk-8 as fiji
USER root

WORKDIR /app

RUN wget https://downloads.openmicroscopy.org/latest/bio-formats5.7/artifacts/bftools.zip
RUN apt-get -y update \
    && apt-get install unzip \
    && unzip bftools.zip -d /opt \
    && rm bftools.zip

RUN mkdir /opt/data /opt/temp

COPY --from=builder /mdc-minerva-image-converter /app

EXPOSE 8080

ENTRYPOINT ["/app/mdc-minerva-image-converter"]
