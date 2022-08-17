package controller

import (
	"testing"

	"github.com/zaaksam/gins/test"
)

func TestOpsTagList(t *testing.T) {
	req := test.NewRequest()

	test.Post(t, req, "/ops/tag/list")
}
