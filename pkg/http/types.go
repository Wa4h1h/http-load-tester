package http

type Header struct {
	Key   string
	Value string
}

type DoResponse struct {
	Url    string
	Status string
	Code   int
	Time   int64
}
