package daemon

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/tunnelshade/rinnegan/agent/log"
)

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", Socket)
}

//HTTPPost lets you just trigger a post call
func HTTPPost(url string, data url.Values) string {
	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}
	log.Debug("Doing POST request to " + url)
	resp, err := client.PostForm("http://d"+url, data)
	if err != nil {
		log.Warn("HTTP call failed")
	}
	var body []byte
	if resp != nil {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	}
	return string(body)
}

//HTTPGet lets you just trigger a post call
func HTTPGet(url string) string {
	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}
	log.Debug("Doing POST request to " + url)
	resp, err := client.Get("http://d" + url)
	if err != nil {
		log.Warn("HTTP call failed")
	}
	var body []byte
	if resp != nil {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	}
	return string(body)
}

//HTTPDelete lets you just trigger a post call
func HTTPDelete(url string) string {
	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}
	log.Debug("Doing DELETE request to " + url)
	req, _ := http.NewRequest("DELETE", "http://d"+url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Warn("HTTP call failed")
	}
	var body []byte
	if resp != nil {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
	}
	return string(body)
}
