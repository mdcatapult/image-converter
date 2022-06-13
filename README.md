# MDC-Minerva-Image-Converter

Webservice to trigger conversion or cropping of images to a valid format for the Minerva UI.

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

Finally, another tool, `bfconvert` is used to convert from a flat .tiff to a pyramidal .ome.tiff 

##To make a pyramid.
bftools in version 6+ has been updated to allow the conversion of flat .tiff images to .ome.tiff pyramid images.

-pyramid-resolutions: This configures how many individual layers the pyramid image file will have (how many tiers to the pyramid)

-pyramid-scale: Defines the reduction in resolution per pyramid resolution

eg -pyramid-resolutions 6 -pyramid-scale 2 of an original image of 10000x5000 pixels would be:

1 - 10000x5000

2 - 5000x2500

3 - 2500x1250

4 - 1250x625

5 - 625x312

6 - 312x156

6 images in which the resolution scales down by a factor of two each time
```
Not to scale depiction of the above pyramid image.

''''''''''''''''''''

  ''''''''''''''''

    ''''''''''''

      ''''''''

       '''''

        '''
        
```
```
bftools/bfconvert -pyramid-resolutions 6 -pyramid-scale 2 pre-pyramid.tiff pyramid.ome.tiff
```
As default the pyramid-resolutions will be set to 6 tiers and the tiers will downscale each time by a factor of 2.

## Deployed Service URL

https://minerva-image-converter.wopr.inf.mdc

In Rancher the deployment is under the `R&D` project in the `minerva` namespace.

## Endpoints

### `/convert`
Converts two seperate .tiff files, one just tissue and the other just roi mask, to a pyramid .ome.tiff file that is 
compatible with the DSP Atlas

Accepts a `Post` request with a JSON body in the format:


```
{
    "input-file" : "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop.tiff",
    "input-mask-file":  "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop-mask.tiff",
    "output-file" : "/opt/data/converted_file.ome.tiff"
}
```

`opt/data` is required as this is specified as the mount point in `image-converter.yml`

### `/crop`
Crops a .tiff image using `bftools` to a given size using specified x and y coords as the center.

Accepts a `Get` request with the following mandatory params:

- x = the x-coordinate to use as the centre of the cropped image. 
- y = the y-coordinate to use as the centre of the cropped image.
- crop-size = the size in pixels to crop the image to.
- experiment-directory = the filepath containing the original raw image, and a channels.pattern file for `bftools` to use to find the image channels.

Returns the cropped image as bytes.

Troubleshooting:

- Out of bounds error: the requested crop falls of the edge of the raw image. x cannot be less than `(cropSize/2)`, and cannot be greater than `({rawImageWidth}-cropSize)`. y is the same but for the height of the raw image rather than its width.

## Testing

### Local

As the main project uses two external conversion tools, the projects dockerfile is built and run via `docker-compose up`

Then tests can be run via:
`go test ./... -v`

### CI 

When running in CI, a custom docker image with docker compose is used as the base image, with a docker in docker service.

The file `docker-compose-ci.yaml` is then used to build and run this project's dockerfile.

This specifies to use the `host` network, which is needed to allow the tests to send requests to the running container.
