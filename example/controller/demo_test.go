package controller

import (
	"testing"

	"github.com/zaaksam/gins/test"
)

func TestDemoList(t *testing.T) {
	req := test.NewRequest()

	test.Post(t, req, "/demo/list")
}
