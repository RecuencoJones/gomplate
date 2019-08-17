//+build integration

package integration

import (
	"net"
	"net/http"

	. "gopkg.in/check.v1"

	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"gotest.tools/v3/icmd"
)

type DatasourcesS3Suite struct {
	l *net.TCPListener
}

var _ = Suite(&DatasourcesS3Suite{})

func (s *DatasourcesS3Suite) SetUpSuite(c *C) {
	backend := s3mem.New()
	s3 := gofakes3.New(backend)
	var err error
	s.l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
	handle(c, err)

	http.Handle("/", s3.Server())
	go http.Serve(s.l, nil)
}

func (s *DatasourcesS3Suite) TearDownSuite(c *C) {
	s.l.Close()
}

func (s *DatasourcesS3Suite) TestReportsVersion(c *C) {
	result := icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/",
		"-H", "foo=Foo:bar",
		"-i", "{{ index (ds `foo`).headers.Foo 0 }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "bar"})

	result = icmd.RunCommand(GomplateBin,
		"-H", "foo=Foo:bar",
		"-i", "{{defineDatasource `foo` `s3://"+s.l.Addr().String()+"/`}}{{ index (ds `foo`).headers.Foo 0 }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "bar"})

	result = icmd.RunCommand(GomplateBin,
		"-i", "{{ $d := ds `s3://"+s.l.Addr().String()+"/`}}{{ index (index $d.headers `Accept-Encoding`) 0 }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "gzip"})
}

func (s *DatasourcesS3Suite) TestTypeOverridePrecedence(c *C) {
	result := icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/foo",
		"-i", "{{ (ds `foo`).value }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})

	result = icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/not.json",
		"-i", "{{ (ds `foo`).value }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "notjson"})

	result = icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/actually.json",
		"-i", "{{ (ds `foo`).value }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})

	result = icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/bogus.csv?type=application/json",
		"-i", "{{ (ds `foo`).value }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})
}

func (s *DatasourcesS3Suite) TestAppendQueryAfterSubPaths(c *C) {
	result := icmd.RunCommand(GomplateBin,
		"-d", "foo=s3://"+s.l.Addr().String()+"/?type=application/json",
		"-i", "{{ (ds `foo` `bogus.csv`).value }}")
	result.Assert(c, icmd.Expected{ExitCode: 0, Out: "json"})
}
