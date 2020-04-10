### Run

`go run brain.go`

### POST data
`curl -X POST --header "TOKEN:<TOKEN>" -d '[{"url":"https://alerty.online","response":200,"request_time":1.343}, {"url":"https://google.com","response":201,"request_time":0.43}]' localhost:3000/monitors/batch`
