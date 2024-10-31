package web_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserHandler Suite")
}
