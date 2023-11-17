package cmds_test

import (
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/cmds"
	"github.com/hsfzxjy/imbed/contrib"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type file struct {
	Path string
	Hash ref.Murmur3
}

func (f file) UploadedPath() string {
	dir, base := path.Split(f.Path)
	return path.Join(dir, f.Hash.WithName(base))
}

type context struct {
	t       *testing.T
	WorkDir string
}

func (c *context) GenImage(seed int64) file {
	im := image.NewRGBA(image.Rect(0, 0, 256, 256))
	rng := rand.New(rand.NewSource(seed))
	rng.Read(im.Pix)
	filename := c.Path(strconv.Itoa(int(seed)) + ".jpg")
	f := util.Unwrap(os.Create(filename))
	defer f.Close()
	util.Check(jpeg.Encode(f, im, &jpeg.Options{Quality: 100}))
	util.Check2(f.Seek(0, 0))
	hash := util.Unwrap(ref.Murmur3FromReader(f))
	return file{filename, hash}
}

type resultTable struct {
	*runResult
	content string
	table   map[string][]string
}

func (t *resultTable) AssertLines(expected int) *resultTable {
	assertLines(t.t, t.content, expected)
	return t
}

func (t *resultTable) Cell(name string, idx int, f func(string)) *resultTable {
	f(t.table[name][idx])
	return t
}

func newResultTable(result *runResult, output string) *resultTable {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	headers := strings.Fields(lines[0])
	table := make(map[string][]string)
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		for i, header := range headers {
			table[header] = append(table[header], fields[i])
		}
	}
	return &resultTable{result, output, table}
}

type runResult struct {
	*context
	Stdout, Stderr string
}

func (r *runResult) Then(f func(r *runResult)) *runResult {
	f(r)
	return r
}

func (r *runResult) Table() *resultTable {
	return newResultTable(r, r.Stdout)
}

func (c *context) Run(args ...string) (result *runResult, err error) {
	registry := transform.NewRegistry()
	contrib.Register(registry)
	commands := app.Commands{}.
		Register(cmds.InitCommand{}.Spec()).
		Register(cmds.AddCommand{}.Spec()).
		Register(cmds.QCommand{}.Spec()).
		Register(cmds.RevCommand{}.Spec())
	args = append(args, "-d", c.WorkDir)
	args = append([]string{"imbed"}, args...)

	oldStdout, oldStderr := app.Stdout, app.Stderr
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}

	defer func() {
		app.Stdout, app.Stderr = oldStdout, oldStderr
		result = &runResult{
			c,
			stdout.String(),
			stderr.String(),
		}
		if err != nil {
			println("-- ERROR --")
			println(err.Error())
		}
		println("-- STDOUT --")
		print(result.Stdout)
		println("-- STDERR --")
		print(result.Stderr)
	}()

	app.Stdout = stdout
	app.Stderr = stderr
	println("!!$", strings.Join(args, " "))
	return nil, app.ParseAndRun(args, commands, registry)
}

func (c *context) RunMust(args ...string) *runResult {
	res, err := c.Run(args...)
	require.ErrorIs(c.t, nil, err)
	return res
}

func (c *context) Path(parts ...string) string {
	parts = append([]string{c.WorkDir}, parts...)
	return path.Join(parts...)
}

func setup(t *testing.T) *context {
	workDir := t.TempDir()

	return &context{WorkDir: workDir, t: t}
}

func Test_Init(t *testing.T) {
	ctx := setup(t)

	ctx.RunMust("init")
	assert.DirExists(t, ctx.Path(".imbed"))
	assert.NoFileExists(t, ctx.Path(".imbed", "db"))
	assert.FileExists(t, ctx.Path("imbed.toml"))
}

func Test_Add(t *testing.T) {
	ctx := setup(t)
	img1 := ctx.GenImage(1)
	var sha string
	var revparsed string
	ctx.
		//
		RunMust("init").
		//
		RunMust("add", img1.Path, "upload.local path =", ctx.Path()).
		Table().AssertLines(2).
		//
		RunMust("q").
		Table().AssertLines(3).
		Cell("SHA", 0, func(s string) {
			assert.Len(t, s, 12)
			sha = s
		}).
		Cell("URL", 1, func(s string) {
			assert.Equal(t, "<none>", s)
		}).
		//
		RunMust("rev", "sha@"+sha).
		Then(func(r *runResult) {
			revparsed = strings.TrimRight(r.Stdout, "\n")
		}).
		//
		RunMust("add", revparsed).
		//
		RunMust("q").Table().AssertLines(3)
	assert.FileExists(t, img1.UploadedPath())
}

func assertLines(t *testing.T, str string, expected int) {
	stripped := strings.TrimRight(str, "\n")
	lines := strings.Split(stripped, "\n")
	assert.Equal(t, expected, len(lines), str)
}
