FROM golang:1.17 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY src ./src
COPY *.go ./

RUN go mod tidy
RUN go build -o /mdc-minerva-image-converter

FROM fiji/fiji:fiji-openjdk-8 as fiji
USER root

WORKDIR /app

RUN wget https://downloads.openmicroscopy.org/bio-formats/6.9.1/artifacts/bftools.zip
RUN apt-get -y update \
    && apt-get install unzip \
    && unzip bftools.zip -d /opt \
    && rm bftools.zip

RUN mkdir /opt/data /opt/temp

COPY --from=builder /mdc-minerva-image-converter /app

EXPOSE 8080

ENV BF_TOOLS_CONVERT_PATH=/opt/bftools/bfconvert
ENV BF_TOOLS_INFO_PATH=/opt/bftools/showinf

ENTRYPOINT ["/app/mdc-minerva-image-converter"]
