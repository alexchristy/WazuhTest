package main

type Rule struct {
	ID          string   `json:"id"`
	Level       int      `json:"level"`
	Description string   `json:"description"`
	Groups      []string `json:"groups,omitempty"`
}

type Output struct {
	Rule       Rule              `json:"rule"`
	Predecoder map[string]string `json:"predecoder,omitempty"`
	Decoder    map[string]string `json:"decoder,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
}

type Data struct {
	Output   Output   `json:"output"`
	Messages []string `json:"messages,omitempty"`
	Token    string   `json:"token,omitempty"`
}

type Response struct {
	Data Data `json:"data"`
}
