package confreaks

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestParseEvents(t *testing.T) {
	b, err := loadFixture("events.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	events, err := ParseEvents(bytes.NewReader(b))
	if err != nil {
		t.Errorf("ParseEvents returned an error %s", err)
	}

	// values based on events.html
	exCount := 151
	count := len(events)
	if count != exCount {
		t.Errorf("There should be %d items, but we found %d items", exCount, count)
	}

	// test last item of events
	event := events[count-1]
	ex := &Event{
		Title:            "MountainWest RubyConf 2007",
		URL:              "http://confreaks.com/events/mwrc2007",
		Date:             time.Date(2007, time.March, 16, 0, 0, 0, 0, time.UTC),
		NumPresentations: 12,
	}

	if !reflect.DeepEqual(event, ex) {
		t.Errorf("%#v is not equal with %#v", event, ex)
	}
}

func TestParseEvent(t *testing.T) {
	b, err := loadFixture("event-details.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	event := &Event{}
	err = ParseEvent(bytes.NewReader(b), event)
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	exCount := 10
	if len(event.Presentations) != exCount {
		t.Errorf("event.Presentations length %d (not %d)", len(event.Presentations), exCount)
	}

	recorded, _ := time.Parse("02-Jan-06 15:04", "19-Oct-13 17:50")
	presentation := &Presentation{
		Title:      "Twitter: Taking Hadoop Realtime with Summingbird",
		Presenters: []string{"Sam Ritchie"},
		URL:        "http://confreaks.com/videos/2768-pnws2013-twitter-taking-hadoop-realtime-with-summingbird",
		Recorded:   recorded,
	}

	lastPresentation := event.Presentations[9]
	if !reflect.DeepEqual(lastPresentation, presentation) {
		t.Errorf("Last presentation %#v is not equal %#v", lastPresentation, presentation)
	}
}

func TestParsePresentation(t *testing.T) {
	b, err := loadFixture("presentation.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	presentation := &Presentation{}
	err = ParsePresentation(bytes.NewReader(b), presentation)
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	ex := &Presentation{
		Title:       "Pharmacist or a Doctor - What does your code base need?",
		VideoURL:    "http://www.youtube.com/watch?v=c-AuPryBSZs",
		Description: "You might know of every single code quality \u0026 metrics tool in the Ruby ecosystem and what they give you, but do you know:\nWhich metrics do you currently need?\nDo you really need them?\nHow do you make your team members own them?\nWait, there was a metaphor in the title\nWhile a pharmacist knows about pretty much every medicine out there and what it cures, its really a doctor who figures out what is required given the symptoms of a patient.\nJust like the vitals recorded for healthy adults, infants, pregnant women or an accident patient in an ICU changes, your code base needs different metrics in different contexts to help you uncover problems.\nTalk take aways\nThrough a lot of examples, the talk helps:\nIdentify the current state of your code base\nUnderstand different metrics and what do they tell you about your code base\nDrive your team towards continuously fixing problems uncovered by these metrics",
	}

	if !reflect.DeepEqual(presentation, ex) {
		t.Errorf("%#v is not equal with %#v", presentation, ex)
	}
}

func TestParsePresentation_html5(t *testing.T) {
	b, err := loadFixture("presentation-html5.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	presentation := &Presentation{}
	err = ParsePresentation(bytes.NewReader(b), presentation)
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	ex := &Presentation{
		Title:       "Your Customers Aren't Stupid and Your Coworkers Are Not Incompetent",
		VideoURL:    "http://cdn.confreaks.com/system/assets/datas/515/original/joe-obrien-small.mp4",
		Description: "Communication is hard. No doubt about it. Many of us, being geeks at heart, have an inherently difficult time communicating with people. Why is it that if we look around, it seems that all we see is incompetence? Why is it that we struggle to get our point across? Why do customers and bosses always seem stupid?\nIn this talk we will focus on communication. How to more effectively listen and speak. We will discuss some patterns that we can look for in ourselves. We will talk about strategies on how can we take a step back and realize what it is we are trying to say and hopefully uncover what it is that our bosses and customers are really hearing.\nI will also walk you through strategies on how to have those difficult conversations and steps that I've learned through my years in sales, consulting, project management and business ownership.",
	}

	if !reflect.DeepEqual(presentation, ex) {
		t.Errorf("%#v is not equal with %#v", presentation, ex)
	}
}

func TestParsePresentation_noVideo(t *testing.T) {
	b, err := loadFixture("presentation-no-video.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	presentation := &Presentation{}
	err = ParsePresentation(bytes.NewReader(b), presentation)
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	ex := &Presentation{
		Title:       "Failing in Plain Sight (Succeeding Invisibly)",
		VideoURL:    "",
		Description: "Supporting every good creative effort is an editorial one, ensuring the\nmessage is focused and clear. This is most successful when nobody notices\nit, and in fact is a failure when people do notice.\nIn the development world, these invisible successes happen in refactoring,\nperformance tuning, spam fighting, security. These areas don't get flashy\nwins, just nasty losses.\nWhen this is what you do, you must get satisfaction from a job well done.\nYou won't get much praise, but what you do get will be sincere, and it\nwill come from people who understand how important it is to keep the",
	}

	if !reflect.DeepEqual(presentation, ex) {
		t.Errorf("%#v is not equal with %#v", presentation, ex)
	}
}
