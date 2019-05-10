// Package gomplate is a template renderer which supports a number of datasources,
// and includes hundreds of built-in functions.
package gomplate

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/hairyhenderson/gomplate/data"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// gomplate -
type gomplate struct {
	funcMap         template.FuncMap
	leftDelim       string
	rightDelim      string
	nestedTemplates templateAliases
	rootTemplate    *template.Template
	context         interface{}
}

// runTemplate -
func (g *gomplate) runTemplate(t *tplate) error {
	tmpl, err := t.toGoTemplate(g)
	if err != nil {
		return err
	}

	// nolint: gocritic
	switch t.target.(type) {
	case io.Closer:
		if t.target != os.Stdout {
			// nolint: errcheck
			defer t.target.(io.Closer).Close()
		}
	}
	err = tmpl.Execute(t.target, g.context)
	return err
}

type templateAliases map[string]string

// newGomplate -
func newGomplate(d *data.Data, leftDelim, rightDelim string, nested templateAliases, context interface{}) *gomplate {
	return &gomplate{
		leftDelim:       leftDelim,
		rightDelim:      rightDelim,
		funcMap:         Funcs(d),
		nestedTemplates: nested,
		context:         context,
	}
}

func parseTemplateArgs(templateArgs []string) (templateAliases, error) {
	nested := templateAliases{}
	for _, templateArg := range templateArgs {
		err := parseTemplateArg(templateArg, nested)
		if err != nil {
			return nil, err
		}
	}
	return nested, nil
}

func parseTemplateArg(templateArg string, ta templateAliases) error {
	parts := strings.SplitN(templateArg, "=", 2)
	pth := parts[0]
	alias := ""
	if len(parts) > 1 {
		alias = parts[0]
		pth = parts[1]
	}

	switch fi, err := fs.Stat(pth); {
	case err != nil:
		return err
	case fi.IsDir():
		files, err := afero.ReadDir(fs, pth)
		if err != nil {
			return err
		}
		prefix := pth
		if alias != "" {
			prefix = alias
		}
		for _, f := range files {
			if !f.IsDir() { // one-level only
				ta[path.Join(prefix, f.Name())] = path.Join(pth, f.Name())
			}
		}
	default:
		if alias != "" {
			ta[alias] = pth
		} else {
			ta[pth] = pth
		}
	}
	return nil
}

// RunTemplates - run all gomplate templates specified by the given configuration
func RunTemplates(o *Config) error {
	Metrics = newMetrics()
	defer runCleanupHooks()
	// make sure config is sane
	o.defaults()
	ds := append(o.DataSources, o.Contexts...)
	d, err := data.NewData(ds, o.DataSourceHeaders)
	if err != nil {
		return err
	}
	addCleanupHook(d.Cleanup)
	nested, err := parseTemplateArgs(o.Templates)
	if err != nil {
		return err
	}
	c, err := createContext(o.Contexts, d)
	if err != nil {
		return err
	}
	g := newGomplate(d, o.LDelim, o.RDelim, nested, c)

	return g.runTemplates(o)
}

func (g *gomplate) runTemplates(o *Config) error {
	start := time.Now()
	tmpl, err := gatherTemplates(o, chooseNamer(o, g))
	Metrics.GatherDuration = time.Since(start)
	if err != nil {
		Metrics.Errors++
		return err
	}
	Metrics.TemplatesGathered = len(tmpl)
	start = time.Now()
	defer func() { Metrics.TotalRenderDuration = time.Since(start) }()
	for _, t := range tmpl {
		tstart := time.Now()
		err := g.runTemplate(t)
		Metrics.RenderDuration[t.name] = time.Since(tstart)
		if err != nil {
			Metrics.Errors++
			return err
		}
		Metrics.TemplatesProcessed++
	}
	return nil
}

func chooseNamer(o *Config, g *gomplate) func(string) (string, error) {
	if o.OutputMap == "" {
		return simpleNamer(o.OutputDir)
	}
	return mappingNamer(o.OutputMap, g)
}

func simpleNamer(outDir string) func(inPath string) (string, error) {
	return func(inPath string) (string, error) {
		outPath := filepath.Join(outDir, inPath)
		return filepath.Clean(outPath), nil
	}
}

func mappingNamer(outMap string, g *gomplate) func(string) (string, error) {
	return func(inPath string) (string, error) {
		out := &bytes.Buffer{}
		t := &tplate{
			name:     "<OutputMap>",
			contents: outMap,
			target:   out,
		}
		tpl, err := t.toGoTemplate(g)
		if err != nil {
			return "", err
		}
		ctx := &context{}
		// nolint: gocritic
		switch c := g.context.(type) {
		case *context:
			for k, v := range *c {
				if k != "in" && k != "ctx" {
					(*ctx)[k] = v
				}
			}
		}
		(*ctx)["ctx"] = g.context
		(*ctx)["in"] = inPath

		err = tpl.Execute(t.target, ctx)
		if err != nil {
			return "", errors.Wrapf(err, "failed to render outputMap with ctx %+v and inPath %s", ctx, inPath)
		}

		return filepath.Clean(strings.TrimSpace(out.String())), nil
	}
}

// RenderTemplate renders a template to the output writer. Data sources and other
// options can be set through 'opts' parameter.
func RenderTemplate(in io.Reader, out io.Writer, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	d, err := data.NewData(nil, nil)
	if err != nil {
		return err
	}
	addCleanupHook(d.Cleanup)
	defer runCleanupHooks()

	i, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	g := newGomplate(d, opts.LDelim, opts.RDelim, nil, nil)

	t := &tplate{
		contents: string(i),
		target:   out,
	}
	err = g.runTemplate(t)
	if err != nil {
		return err
	}
	return nil
}
