package cmds_test

import (
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"path"
	"strconv"
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
	Hash ref.Murmur3Hash
}

func (f file) UploadedPath() string {
	dir, base := path.Split(f.Path)
	return path.Join(dir, ref.FIDFromParts(base, f.Hash).Humanize())
}

type context struct {
	t       *testing.T
	WorkDir string
}

func (c context) GenImage(seed int64) file {
	im := image.NewRGBA(image.Rect(0, 0, 256, 256))
	rng := rand.New(rand.NewSource(seed))
	rng.Read(im.Pix)
	filename := c.Path(strconv.Itoa(int(seed)) + ".jpg")
	f := util.Unwrap(os.Create(filename))
	defer f.Close()
	util.Check(jpeg.Encode(f, im, &jpeg.Options{Quality: 100}))
	util.Check2(f.Seek(0, 0))
	hash := util.Unwrap(ref.Murmur3HashFromReader(f))
	return file{filename, hash}
}

func (c context) Run(args ...string) error {
	registry := transform.NewRegistry()
	contrib.Register(registry)
	commands := app.Commands{}.
		Register(cmds.InitCommand{}.Spec()).
		Register(cmds.AddCommand{}.Spec()).
		Register(cmds.QCommand{}.Spec()).
		Register(cmds.RevCommand{}.Spec())
	args = append(args, "-d", c.WorkDir)
	args = append([]string{"imbed"}, args...)
	return app.ParseAndRun(args, commands, registry)
}

func (c context) RunMust(args ...string) context {
	err := c.Run(args...)
	require.ErrorIs(c.t, nil, err)
	return c
}

func (c context) Path(parts ...string) string {
	parts = append([]string{c.WorkDir}, parts...)
	return path.Join(parts...)
}

func setup(t *testing.T) context {
	workDir := t.TempDir()
	return context{WorkDir: workDir, t: t}
}

func Test_Init(t *testing.T) {
	ctx := setup(t)

	err := ctx.Run("init")
	assert.ErrorIs(t, nil, err)
	assert.DirExists(t, ctx.Path(".imbed"))
	assert.NoFileExists(t, ctx.Path(".imbed", "db"))
	assert.FileExists(t, ctx.Path("imbed.toml"))
}

func Test_Add(t *testing.T) {
	ctx := setup(t)
	img1 := ctx.GenImage(1)
	err := ctx.
		RunMust("init").
		Run("add", img1.Path, "upload.local path =", ctx.Path())
	assert.ErrorIs(t, nil, err)
	assert.FileExists(t, img1.UploadedPath())
}
