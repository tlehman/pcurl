# Results for 2023-05-03 test

## Background
During the first half of the project, I wrote the boilerplate code for making a containerd plugin: https://github.com/tlehman/parpull, before getting the plugin finished though, I realized I should get some empirical data to support the hypothesis that parallelizing the downloads of the layers would actually deliver significant performance gains.

I wrote [pget](https://github.com/tlehman/pget) in Go to test this out using goroutines, but couldn't detect any significant increase by increasing the concurrency. Following the approach that the GSUtil tool uses when `sliced_object_download_max_components` is increased, I wrote a script to use `curl` to do the download (using [HTTP Range headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range)) and then used the bash `&` operator to fork the `curl` off as a background process. Then at the bottom of the script I run `wait` to let the background processes complete the download. In order to focus on **downloading** and not just writing, I used a memory-mapped filesystem, on macos I used this janky one-liner:

```
diskutil eraseVolume HFS+ RAMDisk `hdiutil attach -nomount ram://16384000` 
```

and on Linux (in GCP) I used `/dev/shm`

The [curltest.sh](curltest.sh) script can be run using `./curltest.sh $URL $concurrency` where `$concurrency` is the number of processes you want to fork off. **When running, this script will give a proof-of-concept of how fast a file can be downloaded from storage.googleapis.com to a GCP node, directly into memory**.

## Results
When testing out curltest.sh inside GCP we observed the following:

(All of these download the same 13GiB file)

- Running `time ./curltest.sh 1` downloaded in 64 seconds
- Running `time ./curltest.sh 8` downloaded in 13 seconds (~4x faster)
- Running `time ./curltest.sh 16` downloaded in 11 seconds (~4x faster) 
   - This is likely because the node we tested on only has 8 cores.

When downloading to the disk, there's no speedup, because this node doesn't have an SSD.

## Open questions
1. How does `curltest.sh 8` (setting ofname=/fs/path/to/ssd) perform compared to `./curltest.sh 1` on a GCP node with SSD available?
1. Does translating this approach to Go provide the same benefits?
   - Need to use a forking library (or "os/exec" and a separate binary), but otherwise, if the answer is yes, we can build that into [parpull](https://github.com/tlehman/parpull) and then plug it into containerd to speed up large layer downloads on image pull.
1. How hard is it to configure GKE to use containerd plugins?
   - We can do it manually by running the containerd binary on a GCP node, then by building and running `parpull` and configuring the local containerd to use the parpull socket for gRPC requests. This is how containerd calls out to plugins to perform operations like snapshotting or image downloading.
1. Is it feasible to download _some_ of the layers in parallel directly into memory first?
    - Suppose we have $L$ layers to download, we take take the largest layer, then download directly to memory while the other $L-1$ layers are streaming to the SSD
    - The answer would depend on how much memory is available on the node versus how big the biggest layer is.
1. Assuming 4 is feasible, then, the final question is: how much faster does an image download get when the largest layer is downloaded in parallel directly into memory?
    - This question is inspired by seeing some image downloads that are dominated by one large layer. If there's one layer that's 4x bigger than then average layer, then it's going to put an upper bound on the download speed that's proportional to the largest layer. **But if we parallel-download that largest layer 4x faster than default, then we bring the upper bound down by at most 4x, until we run into another bottle neck (like the second biggest layer, or something else)**