# `pget` is a parallel file downloader

## usage

To build, run `make` and then the `pget` binary will be placed in your current directory

```shell
file_4GB=https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt 
./pget $file_4GB 1
```

## performance

