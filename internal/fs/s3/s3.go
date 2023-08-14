package s3

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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

	s3 := &S3{client: client, bucket: bucket, index: make(map[string]*FileInfo)}
	s3.Index()
	for k := range s3.index {
		log.Printf("%s", k)
	}
	return s3, nil
}

func (s3 *S3) Index() {
	objs := s3.client.ListObjects(context.Background(), s3.bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    "",
	})

	for obj := range objs {
		fpath := path.Join("/", obj.Key)
		log.Printf("%s", obj.Key)
		s3.index[fpath] = &FileInfo{
			name:  filepath.Base(strings.TrimPrefix(fpath, "/")),
			size:  obj.Size,
			isDir: false,
		}

		// Recurse object path and create directories
		dirTree := "/"
		for _, node := range strings.Split(filepath.Dir(fpath), "/") {
			dirTree = filepath.Join(dirTree, node)
			log.Print(dirTree)
			s3.index[dirTree] = &FileInfo{
				name:  filepath.Base(strings.TrimPrefix(dirTree, "/")),
				size:  4096,
				isDir: true,
			}

		}
	}

}
