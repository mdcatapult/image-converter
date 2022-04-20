# MDC-Minerva-Image-Converter

Webservice to trigger conversion of images to a valid format for the Minerva UI.

Currently, a valid image comprises of 6 channels, 3 of which are used by the main 'base' image, with the other 3 being used by the ROI mask image.

This is to allow visually toggling components of the ROI mask on and off, in the Minerva software.

An image processing tool, Fiji is used in a headless mode, which allows a macro file to be used to process the main and mask images to produce one output image.

The macro steps are as follows:

1) Open the main image
2) Split the RGB channels of the main image
3) Open the mask image
4) Split the channels of the mask image
5) Merge the channels of both images, specifying channels 1-3 as the main image, and channels 4-6 as the mask image.
6) Convert to 16 bit, to be compatible with minerva
7) Save

Finally, another tool, `bfconvert` is used to convert from a .tiff to a .ome.tiff 

## Endpoints

### `/convert`
Converts a .tiff file to an .ome.tiff file

Accepts a `Post` request with a JSON body in the format:


```
{
    "input-file" : "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop.tiff",
    "input-mask-file":  "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop-mask.tiff",
    "output-file" : "/opt/data/converted_file.ome.tiff"
}
```

`opt/data` is required as this is specified as the mount point in `image-converter.yml`

## Testing

### Local

As the main project uses two external conversion tools, the projects dockerfile is built and run via `docker-compose up`

Then tests can be run via:
`go test ./... -v`

### CI 

When running in CI, a custom docker image with docker compose is used as the base image, with a docker in docker service.

The file `docker-compose-ci.yaml` is then used to build and run this project's dockerfile.

This specifies to use the `host` network, which is needed to allow the tests to send requests to the running container.
