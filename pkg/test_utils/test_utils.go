package test_utils

import (
	"io/ioutil"
)

const (
	TestVideoPath = "../../../test_data/events/Event_2020-12-27T10:25:05__104__Testing__01.mp4"
	TestImagePath = "../../../test_data/events/Event_2020-12-27T10:25:09__104__Testing__01.jpg"
)

func GetTempDir() (string, error) {
	return ioutil.TempDir("", "cameranator")
}
