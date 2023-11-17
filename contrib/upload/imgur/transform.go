package imgur

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen upload.imgur#Applier
type Applier struct {
	App `imbed:""`
}

const API = "https://api.imgur.com/3/image"

var apiUrl, _ = url.Parse(API)

func (u *Applier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	client, err := util.ClientWithProxy(app.ProxyFunc(), apiUrl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	formData := multipart.NewWriter(&buf)

	err = formData.WriteField("type", "base64")
	if err != nil {
		return nil, err
	}

	fhash, err := a.Content().GetHash()
	if err != nil {
		return nil, err
	}
	filename := fhash.WithName(a.BaseName())
	err = formData.WriteField("name", filename)
	if err != nil {
		return nil, err
	}

	{
		w, err := formData.CreateFormField("image")
		if err != nil {
			return nil, err
		}
		b64e := base64.NewEncoder(base64.RawStdEncoding, w)
		r, err := a.Content().BytesReader()
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(b64e, r)
		if err != nil {
			return nil, err
		}
		err = b64e.Close()
		if err != nil {
			return nil, err
		}
	}

	err = formData.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", API, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", formData.FormDataContentType())
	req.Header.Add("Host", "api.imgur.com")
	req.Header.Add("User-Agent", "Imbed")
	req.Header.Add("Authorization", "Client-ID "+u.ClientId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBodyRaw, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("imgur: bad response status=%s, body=%s", resp.Status, string(respBodyRaw))
	}

	jsonD := json.NewDecoder(bytes.NewReader(respBodyRaw))
	var respBody struct {
		Data struct {
			Link       string
			DeleteHash string
		}
	}
	err = jsonD.Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return asset.MergeUpdates(
		asset.UpdateExt([]byte(respBody.Data.DeleteHash)),
		asset.UpdateUrl(respBody.Data.Link),
	), nil
}

//imbed:schemagen
type App struct {
	ClientId string `imbed:"clientId"`
}

//imbed:schemagen upload.imgur#Config
type Config struct {
	Apps    map[string]App `imbed:"apps"`
	Default string         `imbed:"default,\"\""`
}

func (c *Config) Validate() error {
	n := len(c.Apps)
	if n == 0 {
		return errors.New("imgur: no available app")
	}
	if c.Default != "" {
		_, ok := c.Apps[c.Default]
		if !ok {
			return fmt.Errorf("imgur: bad default app name %q", c.Default)
		}
	} else if n == 1 {
		for k := range c.Apps {
			c.Default = k
			break
		}
	} else {
		return fmt.Errorf("imgur: no default app specified")
	}
	return nil
}

//imbed:schemagen upload.imgur#Params
type Params struct {
	AppName string `imbed:"app,\"\""`
}

func (p *Params) BuildTransform(c *Config) (transform.Applier, error) {
	appName := p.AppName
	if appName == "" {
		appName = c.Default
	}
	app, ok := c.Apps[appName]
	if !ok {
		return nil, fmt.Errorf("imgur: no app named %q", appName)
	}
	return &Applier{App: app}, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "upload.imgur",
		ConfigSchema.Build(),
		ParamsSchema.Build(),
	).
		Alias("imgur").
		Category(transform.Terminal)
}
