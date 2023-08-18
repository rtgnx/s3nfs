package s3

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	Alias, Endpoint, SecretKey, AccessKey, Bucket string
}

func (cfg S3Config) Open() (billy.Filesystem, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: true,
	})

	if err != nil {
		return nil, err
	}

	s3 := &S3{client: client, bucket: cfg.Bucket, index: make(map[string]fs.FileInfo)}
	s3.Index()
	for k := range s3.index {
		log.Printf("%s", k)
	}
	return s3, nil
}

func FromEnv() (billy.Filesystem, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY")
	secretAccessKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")
	useSSL := true

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	s3 := &S3{client: client, bucket: bucket, index: make(map[string]fs.FileInfo)}
	s3.Index()
	for k := range s3.index {
		log.Printf("%s", k)
	}
	return s3, nil
}

func (s3 *S3) Index() {

	go func() {

		objs := s3.client.ListObjects(context.Background(), s3.bucket, minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    "",
		})

		for obj := range objs {
			fpath := path.Join("/", obj.Key)
			log.Printf("%s", obj.Key)
			s3.index[fpath] = NewFileInfo(obj)

			// Recurse object path and create directories
			dirTree := "/"
			for _, node := range strings.Split(filepath.Dir(fpath), "/") {
				dirTree = filepath.Join(dirTree, node)
				log.Print(dirTree)
				s3.index[dirTree] = NewDirInfo(dirTree)

			}
		}
		<-time.After(time.Second * 60)
	}()

}
