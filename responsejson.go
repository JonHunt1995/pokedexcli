package main

type result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type responseJSON struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []result `json:"results"`
}
