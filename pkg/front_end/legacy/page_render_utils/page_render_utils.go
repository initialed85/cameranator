package page_render_utils

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/initialed85/cameranator/pkg/persistence/legacy"
)

type SummaryTableRowSeed struct {
	EventsURL  string
	EventsDate string
	EventCount string
}

func renderSummaryTableRows(eventsByDate map[time.Time][]legacy.Event) (string, error) {
	keys := make([]time.Time, 0)
	for key := range eventsByDate {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Unix() > keys[j].Unix()
	})

	b := bytes.Buffer{}
	for _, eventsDate := range keys {
		events := eventsByDate[eventsDate]

		t := template.New("SummaryTableRowSeed")

		t, err := t.Parse(SummaryTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsSummaryTableRowSeed := SummaryTableRowSeed{
			EventsURL:  fmt.Sprintf("events_%v.html", eventsDate.Format("2006_01_02")),
			EventsDate: eventsDate.Format("2006-01-02"),
			EventCount: fmt.Sprintf("%v", len(events)),
		}

		err = t.Execute(&b, eventsSummaryTableRowSeed)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(b.String()), nil
}

type SummarySeed struct {
	Title      string
	Now        string
	StyleSheet string
	TableRows  string
}

func RenderSummary(title string, eventsByDate map[time.Time][]legacy.Event, now time.Time) (string, error) {
	t := template.New("SummarySeed")

	t, err := t.Parse(SummaryHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	tableRows, err := renderSummaryTableRows(eventsByDate)
	if err != nil {
		return "", err
	}

	eventSummary := SummarySeed{
		Title:      title,
		Now:        now.Format("2006-01-02 15:04:05"),
		StyleSheet: StyleSheet,
		TableRows:  tableRows,
	}

	err = t.Execute(&b, eventSummary)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(b.String()), nil
}

type CheckBoxSeed struct {
	CameraName string
}

func renderCheckBoxes(events []legacy.Event) (string, error) {
	cameraNamesMap := make(map[string]bool)
	for _, event := range events {
		cameraNamesMap[event.CameraName] = true
	}

	cameraNames := make([]string, 0)
	for cameraName := range cameraNamesMap {
		cameraNames = append(cameraNames, cameraName)
	}

	sort.Strings(cameraNames)

	b := bytes.Buffer{}
	for _, cameraName := range cameraNames {
		t := template.New("CheckBoxSeed")

		t, err := t.Parse(CheckBoxHTML)
		if err != nil {
			return "", err
		}

		checkBoxSeed := CheckBoxSeed{
			CameraName: cameraName,
		}

		err = t.Execute(&b, checkBoxSeed)
		if err != nil {
			return "", nil
		}

		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String()), nil
}

type PageTableRowSeed struct {
	EventID         string
	Timestamp       string
	Size            string
	CameraName      string
	HighResImageURL string
	LowResImageURL  string
	HighResVideoURL string
	LowResVideoURL  string
}

func renderPageTableRows(events []legacy.Event) (string, error) {
	b := bytes.Buffer{}
	for _, event := range events {
		t := template.New("PageTableRowSeed")

		t, err := t.Parse(PageTableRowHTML)
		if err != nil {
			return "", err
		}

		eventsTableRowSeed := PageTableRowSeed{
			EventID:         event.EventID.String(),
			Timestamp:       event.Timestamp.Format("2006-01-02 15:04:05"),
			Size:            "?",
			CameraName:      event.CameraName,
			HighResImageURL: event.HighResImagePath,
			LowResImageURL:  event.LowResImagePath,
			HighResVideoURL: event.HighResVideoPath,
			LowResVideoURL:  event.LowResVideoPath,
		}

		err = t.Execute(&b, eventsTableRowSeed)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimRight(b.String(), " \r\n\t"), nil
}

type PageSeed struct {
	Title      string
	EventsDate string
	Now        string
	StyleSheet string
	JavaScript string
	CheckBoxes string
	TableRows  string
}

func RenderPage(title string, events []legacy.Event, eventsDate, now time.Time) (string, error) {
	t := template.New("PageSeed")

	t, err := t.Parse(PageHTML)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}

	checkBoxes, err := renderCheckBoxes(events)
	if err != nil {
		return "", err
	}

	tableRows, err := renderPageTableRows(events)
	if err != nil {
		return "", err
	}

	eventsSeed := PageSeed{
		Title:      title,
		EventsDate: eventsDate.Format("2006-01-02"),
		Now:        now.Format("2006-01-02 15:04:05"),
		StyleSheet: StyleSheet,
		JavaScript: JavaScript,
		CheckBoxes: checkBoxes,
		TableRows:  tableRows,
	}

	err = t.Execute(&b, eventsSeed)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(b.String()), nil
}
