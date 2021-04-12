package main

type Stock struct {
	Ticker string  `json:"ticker"`
	Title  string  `json:"title"`
	Rsi    float32 `json:"rsi"`
}
