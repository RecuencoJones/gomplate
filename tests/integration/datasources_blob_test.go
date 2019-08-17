//+build integration

package integration

import (
	. "gopkg.in/check.v1"

	"gotest.tools/v3/icmd"
)

type BlobDatasourcesSuite struct {
	// l *net.TCPListener
}

var _ = Suite(&BlobDatasourcesSuite{})

// func (s *BlobDatasourcesSuite) SetUpSuite(c *C) {
// 	var err error
// 	s.l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
// 	handle(c, err)

// 	http.HandleFunc("/", mirrorHandler)
// 	http.HandleFunc("/not.json", typeHandler("application/yaml", "value: notjson\n"))
// 	http.HandleFunc("/foo", typeHandler("application/json", `{"value": "json"}`))
// 	http.HandleFunc("/actually.json", typeHandler("", `{"value": "json"}`))
// 	http.HandleFunc("/bogus.csv", typeHandler("text/plain", `{"value": "json"}`))
// 	go http.Serve(s.l, nil)
// }

// func (s *BlobDatasourcesSuite) TearDownSuite(c *C) {
// 	s.l.Close()
// }

func (s *BlobDatasourcesSuite) TestS3Datasource(c *C) {
	result := icmd.RunCommand(GomplateBin,
		"-d", "passengers=s3://ryft-public-sample-data/passengers.txt?region=us-east-1",
		"-i", "{{ include `passengers` }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "bar"})

	// result = icmd.RunCommand(GomplateBin,
	// 	"-H", "foo=Foo:bar",
	// 	"-i", "{{defineDatasource `foo` `http://"+s.l.Addr().String()+"/`}}{{ index (ds `foo`).headers.Foo 0 }}")
	// result.Assert(c, icmd.Expected{ExitCode: 0, Out: "bar"})

	// result = icmd.RunCommand(GomplateBin,
	// 	"-i", "{{ $d := ds `http://"+s.l.Addr().String()+"/`}}{{ index (index $d.headers `Accept-Encoding`) 0 }}")
	// result.Assert(c, icmd.Expected{ExitCode: 0, Out: "gzip"})
}

// func (s *BlobDatasourcesSuite) TestTypeOverridePrecedence(c *C) {
// 	result := icmd.RunCommand(GomplateBin,
// 		"-d", "foo=http://"+s.l.Addr().String()+"/foo",
// 		"-i", "{{ (ds `foo`).value }}")
// 	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})

// 	result = icmd.RunCommand(GomplateBin,
// 		"-d", "foo=http://"+s.l.Addr().String()+"/not.json",
// 		"-i", "{{ (ds `foo`).value }}")
// 	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "notjson"})

// 	result = icmd.RunCommand(GomplateBin,
// 		"-d", "foo=http://"+s.l.Addr().String()+"/actually.json",
// 		"-i", "{{ (ds `foo`).value }}")
// 	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})

// 	result = icmd.RunCommand(GomplateBin,
// 		"-d", "foo=http://"+s.l.Addr().String()+"/bogus.csv?type=application/json",
// 		"-i", "{{ (ds `foo`).value }}")
// 	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})
// }

// func (s *BlobDatasourcesSuite) TestAppendQueryAfterSubPaths(c *C) {
// 	result := icmd.RunCommand(GomplateBin,
// 		"-d", "foo=http://"+s.l.Addr().String()+"/?type=application/json",
// 		"-i", "{{ (ds `foo` `bogus.csv`).value }}")
// 	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})
// }
