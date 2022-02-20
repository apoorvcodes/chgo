package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getLessonsWithSource(s string) (lessons string, source string) {
	ts := strings.TrimSpace(s)
	tokens := strings.Split(ts, "\n")

	var processed []string
	for _, token := range tokens {
		t := strings.TrimSpace(token)

		if len(t) != 0 {
			processed = append(processed, t)
		}
	}

	if len(processed) == 1 {
		source = processed[0]
		return
	}

	return processed[0], processed[1]
}

func ExtractCourses(r io.Reader) ([]Course, error) {
	var courses []Course

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("article.course").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a.course-btn").Attr("href")

		title := s.Find("h3").Text()
		title = strings.TrimSpace(title)

		duration := s.Find(".course-duration").Text()
		duration = strings.TrimSpace(duration)

		language := s.Find(".course-lang").Text()
		language = strings.TrimSpace(language)

		ls := s.Find(".course-lessons").Text()
		lessons, source := getLessonsWithSource(ls)

		courses = append(courses, Course{
			Title:    title,
			URL:      url,
			Duration: duration,
			Language: language,
			Lessons:  lessons,
			Source:   source,
		})
	})

	return courses, nil
}

func extractCourseID(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", nil
	}

	courseID, ok := doc.Find("button.course-action").Attr("data-course-id")
	if !ok {
		return "", fmt.Errorf("failed to extract course id")
	}

	return courseID, err
}
