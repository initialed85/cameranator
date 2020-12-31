package page_render_utils

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/initialed85/cameranator/pkg/persistence/legacy"
)

func getEventsAndEventsDescendingByDateDescendingAndLastTimeAndNow() ([]legacy.Event, map[time.Time][]legacy.Event, time.Time, time.Time) {
	time1 := time.Time{}
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)
	time4 := time1.Add(time.Minute * 24)
	time5 := time2.Add(time.Minute * 24)
	time6 := time3.Add(time.Minute * 24)

	event1 := legacy.Event{Timestamp: time1, CameraName: "camera1", HighResImagePath: "image1-hi", LowResImagePath: "image1-lo", HighResVideoPath: "video1-hi", LowResVideoPath: "video1-lo"}
	event2 := legacy.Event{Timestamp: time2, CameraName: "camera2", HighResImagePath: "image2-hi", LowResImagePath: "image2-lo", HighResVideoPath: "video2-hi", LowResVideoPath: "video2-lo"}
	event3 := legacy.Event{Timestamp: time3, CameraName: "camera3", HighResImagePath: "image3-hi", LowResImagePath: "image3-lo", HighResVideoPath: "video3-hi", LowResVideoPath: "video3-lo"}
	event4 := legacy.Event{Timestamp: time4, CameraName: "camera4", HighResImagePath: "image4-hi", LowResImagePath: "image4-lo", HighResVideoPath: "video4-hi", LowResVideoPath: "video4-lo"}
	event5 := legacy.Event{Timestamp: time5, CameraName: "camera5", HighResImagePath: "image5-hi", LowResImagePath: "image5-lo", HighResVideoPath: "video5-hi", LowResVideoPath: "video5-lo"}
	event6 := legacy.Event{Timestamp: time6, CameraName: "camera6", HighResImagePath: "image6-hi", LowResImagePath: "image6-lo", HighResVideoPath: "video6-hi", LowResVideoPath: "video6-lo"}

	event1.EventID = uuid.UUID{0}
	event2.EventID = uuid.UUID{1}
	event3.EventID = uuid.UUID{2}
	event4.EventID = uuid.UUID{3}
	event5.EventID = uuid.UUID{4}
	event6.EventID = uuid.UUID{5}

	events := []legacy.Event{
		event1,
		event2,
		event3,
		event4,
		event5,
		event6,
	}

	eventsDescendingByDateDescending := legacy.GetEventsDescendingByDateDescending(events)

	events = legacy.GetEventsDescending(events)

	return events, eventsDescendingByDateDescending, time3, time3.Add(time.Hour * 24)
}

func TestRenderSummary(t *testing.T) {
	_, eventsDescendingByDateDescending, _, now := getEventsAndEventsDescendingByDateDescendingAndLastTimeAndNow()

	data, err := RenderSummary("All events", eventsDescendingByDateDescending, now)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\n" + data + "\n\n")

	assert.Equal(
		t,
		"<html>\n<head>\n<title>All events as at 0001-01-04 00:00:00</title>\n\n<style type=\"text/css\">\nBODY {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: normal;\n\ttext-align: center;\n}\n\nTH {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: bold;\n\ttext-align: center;\n}\n\nTD {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: normal;\n\ttext-align: center;\n\tborder: 1px solid gray; \n}\n</style>\n</head>\n\n<body>\n<h2>All events as at 0001-01-04 00:00:00</h2>\n\n<center>\n\n<table width=\"90%\">\n\n\t<tr>\n\t\t<th>Date</th>\n\t\t<th>Events</th>\n\t</tr>\n<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_03.html\">0001-01-03</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_02.html\">0001-01-02</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_01.html\">0001-01-01</a></td>\n\t\t<td>2</td>\n\t</tr>\n</table>\n\n</center>\n</body>\n</html>",
		data,
	)
}

func TestRenderPage(t *testing.T) {
	_, eventsDescendingByDateDescending, time3, now := getEventsAndEventsDescendingByDateDescendingAndLastTimeAndNow()

	data, err := RenderPage("Events", eventsDescendingByDateDescending[time3], time3, now)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\n" + data + "\n\n")

	assert.Equal(
		t,
		"<html>\n<head>\n<title>Events for 0001-01-03 as at 0001-01-04 00:00:00</title>\n\n<style type=\"text/css\">\nBODY {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: normal;\n\ttext-align: center;\n}\n\nTH {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: bold;\n\ttext-align: center;\n}\n\nTD {\n\tfont-family: Tahoma, serif;\n\tfont-size: 8pt;\n\tfont-weight: normal;\n\ttext-align: center;\n\tborder: 1px solid gray; \n}\n</style>\n\n<script>\nfunction toggleCamera(cameraName) {\n\tArray.from(document.getElementsByClassName(cameraName)).map((x) => {\n\t\tx.style.display = x.style.display === 'none' ? '' : 'none'\n\t})\n}\n</script>\n</head>\n\n<body>\n<h1>Events for 0001-01-03 as at 0001-01-04 00:00:00</h1>\n\n<center>\n\ncamera3 <input type=\"checkbox\" checked=\"true\" onclick=\"toggleCamera('camera3')\"/>\ncamera6 <input type=\"checkbox\" checked=\"true\" onclick=\"toggleCamera('camera6')\"/>\n\n<br />\n<br />\n\n<table width=\"90%\">\n\t<tr>\n\t\t<th>Event ID</th>\n\t\t<th>Timestamp</th>\n\t\t<th>Camera</th>\n\t\t<th>Screenshot</th>\n\t\t<th>Download</th>\n\t</tr>\n\n\t<tr class=\"camera6\">\n\t\t<td>05000000-0000-0000-0000-000000000000</td>\n\t\t<td>0001-01-03 00:24:00</td>\n\t\t<td>camera6</td>\n\t\t<td style=\"width: 320px\";><a target=\"_blank\" href=\"image6-hi\"><img src=\"image6-lo\" width=\"320\" height=\"180\" /></a></td>\n\t\t<td>Download <a target=\"_blank\" href=\"video6-hi\">high-res</a> or <a target=\"_blank\" href=\"video6-lo\">low-res</a></td>\n\t</tr>\n\n\t<tr class=\"camera3\">\n\t\t<td>02000000-0000-0000-0000-000000000000</td>\n\t\t<td>0001-01-03 00:00:00</td>\n\t\t<td>camera3</td>\n\t\t<td style=\"width: 320px\";><a target=\"_blank\" href=\"image3-hi\"><img src=\"image3-lo\" width=\"320\" height=\"180\" /></a></td>\n\t\t<td>Download <a target=\"_blank\" href=\"video3-hi\">high-res</a> or <a target=\"_blank\" href=\"video3-lo\">low-res</a></td>\n\t</tr>\n</table>\n\n<center>\n</body>\n</html>",
		data,
	)
}
