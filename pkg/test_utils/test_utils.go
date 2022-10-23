package test_utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	TestVideoPath = "../../../test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4"
	TestImagePath = "../../../test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg"
)

func GetTempDir() (string, error) {
	return ioutil.TempDir("", "cameranator")
}

func IsLive(host string, port int64) bool {
	resp, err := http.Get(
		fmt.Sprintf(
			"http://%v:%v/healthz",
			host,
			port,
		),
	)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
