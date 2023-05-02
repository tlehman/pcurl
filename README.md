# `pget` is a parallel file downloader

## usage

To build, run `make` and then the `pget` binary will be placed in your current directory

```shell
file_4GB=https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt 
./pget $file_4GB 1
```

## performance

```shell
% ./pget http://defini.dev.s3-website-us-west-2.amazonaws.com/sd-v1-4.ckpt 1
File downloaded in 51 seconds at 79 MiB/s!
% ./pget http://defini.dev.s3-website-us-west-2.amazonaws.com/sd-v1-4.ckpt 8
File downloaded in 49 seconds at 83 MiB/s!
% ./pget http://defini.dev.s3-website-us-west-2.amazonaws.com/sd-v1-4.ckpt 32
File downloaded in 44 seconds at 92 MiB/s!
```

Needs more profiling, this is not _that_ signifcant of an improvement