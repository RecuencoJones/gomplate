package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadBlob(t *testing.T) {
	os.Setenv("AWS_ANON", "true")
	d, err := NewData([]string{"-d", "passengers=s3://ryft-public-sample-data/AWS-x86-AMI-queries.json?region=us-east-1&type=application/array%2Bjson"}, nil)
	assert.NoError(t, err)

	out, err := d.Datasource("passengers")
	assert.NoError(t, err)
	assert.NotNil(t, out)
}
