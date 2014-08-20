package confreaks

type byDate []*Event

func (b byDate) Len() int           { return len(b) }
func (b byDate) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byDate) Less(i, j int) bool { return b[i].Date.After(b[j].Date) }

type byRecorded []*Presentation

func (b byRecorded) Len() int           { return len(b) }
func (b byRecorded) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byRecorded) Less(i, j int) bool { return b[i].Recorded.Before(b[j].Recorded) }
