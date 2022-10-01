module example.com/example

go 1.19

require (
	example.com/muxjwt v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)

replace example.com/muxjwt => ../src
