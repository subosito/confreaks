package confreaks

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPresentation_ParseDetails(t *testing.T) {
	b, err := loadFixture("presentation.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	presentation := &Presentation{}
	err = presentation.ParseDetails(bytes.NewReader(b))
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

func TestPresentation_ParseDetailsHTML5(t *testing.T) {
	b, err := loadFixture("presentation-html5.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	presentation := &Presentation{}
	err = presentation.ParseDetails(bytes.NewReader(b))
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
