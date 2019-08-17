package data

import (
	"context"
	"path"

	"github.com/aws/aws-sdk-go/aws"

	gaws "github.com/hairyhenderson/gomplate/aws"
	"github.com/pkg/errors"

	"gocloud.dev/blob/s3blob"
)

func parseBlobArgs(origPath string, args ...string) (paramPath string, err error) {
	paramPath = origPath
	if len(args) >= 1 {
		paramPath = path.Join(paramPath, args[0])
	}

	if len(args) >= 2 {
		err = errors.New("Maximum two arguments to s3 datasource: alias, extraPath")
	}
	return paramPath, err
}

func readBlob(source *Source, args ...string) (output []byte, err error) {
	ctx := context.TODO()

	paramPath, err := parseBlobArgs(source.URL.Path, args...)
	if err != nil {
		return nil, err
	}

	sess := gaws.SDKSession()
	r := source.URL.Query().Get("region")
	if r != "" {
		sess.Config.Region = aws.String(r)
	}

	b := source.URL.Host

	bucket, err := s3blob.OpenBucket(ctx, sess, b, nil)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	return bucket.ReadAll(ctx, paramPath)
}
