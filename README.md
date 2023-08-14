# NFS bridge to S3

Read-only nfs server that acts as a bridge to S3 bucket

##  Install 

```
go install github.com/rtgnx/s3nfs/cmd/s3nfsd@latest

```



## Usage


```
S3_ENDPOINT=s3.fr-par.scw.cloud
S3_ACCESS_KEY=<ACCESS_KEY>
S3_SECRET_KEY=<SECRET_KEY>
S3_BUCKET=<BUCKET>

s3nfsd serve --addr :8888
```


```
mkdir -p $HOME/s3
mount -o port=8888,mountport=8888 -t nfs 127.0.0.1:/ $HOME/s3
```