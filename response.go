package gins

const (
	gins_context_response_api = "gins_context_response_api"
	gins_context_response_web = "gins_context_response_web"
)

type IResponse interface {
	render()
}
