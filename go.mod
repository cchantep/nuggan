module github.com/cchantep/nuggan

go 1.13

replace nuggan => ./src

require (
	github.com/aws/aws-lambda-go v1.13.3 // indirect
	github.com/davidbyttow/govips v0.0.0-20190304175058-d272f04c0fea
	github.com/pelletier/go-toml v1.6.0 // indirect
	nuggan v0.0.0-00010101000000-000000000000
)
