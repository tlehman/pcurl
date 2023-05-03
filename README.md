# `pget` is a parallel file downloader

## usage

To build, run `make` and then the `pget` binary will be placed in your current directory

```shell
file_4GB=https://huggingface.co/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt 
./pget $file_4GB 1
```

# curltest (for comparison)

Running `./curltest.sh 16` spawns 16 processes, each downloading a different section of the file. Runs consistently around 45s (min 44.1, max 45.9).

Running `./curltest.sh 1` spawns 1 process. Runs consistently around 56s (min 53.4, max 59.1).

The `curltest.sh` script is meant to quickly test out the performance of parallelizing a file download by using the HTTP Range header and letting the OS handle the process scheduling by forking `concurrency` number of times.

```
% time ./curltest.sh 16
./curltest.sh 16  2.10s user 10.20s system 27% cpu 44.107 total
% time ./curltest.sh 16
./curltest.sh 16  1.86s user 9.23s system 24% cpu 45.978 total
% time ./curltest.sh 16
./curltest.sh 16  2.62s user 12.63s system 33% cpu 44.924 total
% time ./curltest.sh 16
./curltest.sh 16  2.50s user 12.00s system 32% cpu 44.706 total

% time ./curltest.sh 1
./curltest.sh 1  1.05s user 5.49s system 11% cpu 56.614 total
% time ./curltest.sh 1
./curltest.sh 1  1.06s user 5.99s system 12% cpu 55.621 total
% time ./curltest.sh 1
./curltest.sh 1  1.10s user 6.68s system 13% cpu 59.088 total
% time ./curltest.sh 1
./curltest.sh 1  1.05s user 6.03s system 13% cpu 53.379 total
```

The conclusion is that parallelization from 1 to 16 download processes reduces download time  about 20% (since $1 - 45s/56s = .19$).

[2023-03-05 test results](testresults.md)