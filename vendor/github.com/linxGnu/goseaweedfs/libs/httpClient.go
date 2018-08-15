package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var reg, _ = regexp.Compile("[^a-zA-Z0-9.-_]+")

// MakeURL encode to full url request
func MakeURL(scheme, host, path string, args url.Values) string {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	if args != nil {
		u.RawQuery = args.Encode()
	}
	return u.String()
}

// HTTPClient wrapper for http client
type HTTPClient struct {
	Client *http.Client
}

// NewHTTPClient new http client wrapper
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{Client: &http.Client{
		Timeout: timeout,
	}}
}

// // PostBytes post raw bytes
// func (c *HTTPClient) PostBytes(url string, body []byte) (respBody []byte, statusCode int, err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			err = fmt.Errorf("%v", e)
// 		}
// 	}()

// 	r, err := c.Client.Post(url, "application/octet-stream", bytes.NewReader(body))
// 	if err != nil {
// 		err = fmt.Errorf("Post to %s: %v", url, err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	statusCode = r.StatusCode
// 	respBody, err = ioutil.ReadAll(r.Body)

// 	return
// }

// PostForm do post with post form values
func (c *HTTPClient) PostForm(_url string, values url.Values) (respBody []byte, statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	if values == nil {
		values = make(url.Values)
	}

	r, err := c.Client.PostForm(_url, values)
	if err != nil {
		err = fmt.Errorf("Post to %s: %v", _url, err)
		return
	}
	defer r.Body.Close()

	statusCode = r.StatusCode
	respBody, err = ioutil.ReadAll(r.Body)

	return
}

// Get make get request
func (c *HTTPClient) Get(scheme, host, path string, values url.Values) (respBody []byte, statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	if values == nil {
		values = make(url.Values)
	}

	return c.GetWithURL(MakeURL(scheme, host, path, values))
}

// GetWithHeaders do get with customer headers
func (c *HTTPClient) GetWithHeaders(fullURL string, headers map[string]string) (respBody []byte, statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		err = fmt.Errorf("Get %s: %v", fullURL, err)
		return
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	r, e := c.Client.Do(req)
	if e != nil {
		err = fmt.Errorf("Delete %s: %v", fullURL, e)
		return
	}
	defer r.Body.Close()

	statusCode = r.StatusCode

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		err = fmt.Errorf("Delete %s: %v", fullURL, e)
		return
	}
	respBody = body

	return
}

// GetWithURL do http get with full url/uri
func (c *HTTPClient) GetWithURL(fullURL string) (respBody []byte, statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	r, err := c.Client.Get(fullURL)
	if err != nil {
		err = fmt.Errorf("Get from %s: %v", fullURL, err)
		return
	}
	defer r.Body.Close()

	statusCode = r.StatusCode
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s: %s", fullURL, r.Status)
		return
	}

	respBody, err = ioutil.ReadAll(r.Body)

	return
}

// Delete make delete method request
func (c *HTTPClient) Delete(url string) (statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		err = fmt.Errorf("Delete %s: %v", url, err)
		return
	}

	r, e := c.Client.Do(req)
	if e != nil {
		err = fmt.Errorf("Delete %s: %v", url, e)
		return
	}
	defer r.Body.Close()

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		err = fmt.Errorf("Delete %s: %v", url, e)
		return 0, err
	}

	switch r.StatusCode {
	case http.StatusNotFound, http.StatusAccepted, http.StatusOK:
		err = nil
		return
	}

	m := make(map[string]interface{})
	if e := json.Unmarshal(body, &m); e == nil {
		if s, ok := m["error"].(string); ok {
			err = fmt.Errorf("Delete %s: %v", url, s)
			return
		}
	}

	err = fmt.Errorf("Delete %s. Got response but can not parse.", url)
	return
}

// // DownloadFromURL download file from url.
// // Note: rc must be closed after finishing as other ReadCloser.
// func (c *HTTPClient) DownloadFromURL(fileURL string) (filename string, rc io.ReadCloser, err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			err = fmt.Errorf("%v", e)
// 		}
// 	}()

// 	r, err := c.Client.Get(fileURL)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	if r.StatusCode != http.StatusOK {
// 		r.Body.Close()
// 		err = fmt.Errorf("Download %s: %s", fileURL, r.Status)
// 		return
// 	}

// 	contentDisposition := r.Header["Content-Disposition"]
// 	if len(contentDisposition) > 0 {
// 		if strings.HasPrefix(contentDisposition[0], "filename=") {
// 			filename = contentDisposition[0][len("filename="):]
// 			filename = strings.Trim(filename, "\"")
// 		}
// 	}
// 	rc = r.Body

// 	return
// }

// // Do any request
// func (c *HTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			err = fmt.Errorf("%v", e)
// 		}
// 	}()

// 	return c.Client.Do(req)
// }

// Upload file content
func (c *HTTPClient) Upload(uploadURL string, filename string, reader io.Reader, isGzipped bool, mtype string) (respBody []byte, statusCode int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	return c.uploadContent(uploadURL, func(w io.Writer) (err error) {
		_, err = io.Copy(w, reader)
		return
	}, filename, isGzipped, mtype)
}

func (c *HTTPClient) uploadContent(uploadURL string, fillBuffer func(w io.Writer) error, filename string, isGzipped bool, mtype string) (respBody []byte, statusCode int, err error) {
	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, reg.ReplaceAllString(filename, "")))
	if mtype == "" {
		mtype = mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))
	}
	if mtype != "" {
		h.Set("Content-Type", mtype)
	}
	if isGzipped {
		h.Set("Content-Encoding", "gzip")
	}

	fileWriter, err := bodyWriter.CreatePart(h)
	if err != nil {
		return
	} else if err = fillBuffer(fileWriter); err != nil {
		return
	}

	contentType := bodyWriter.FormDataContentType()
	if err = bodyWriter.Close(); err != nil {
		return
	}

	r, err := c.Client.Post(uploadURL, contentType, body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	statusCode = r.StatusCode
	respBody, err = ioutil.ReadAll(r.Body)

	return
}
