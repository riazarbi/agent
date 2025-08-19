module agent

go 1.22.1

require (
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.3.0
	github.com/chzyer/readline v1.5.1
	github.com/invopop/jsonschema v0.13.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/JohannesKaufmann/dom v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/openai/openai-go/v2 v2.0.2
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace agent/internal/editcorrector => ./internal/editcorrector
