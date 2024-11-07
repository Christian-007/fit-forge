package services_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTodoService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TodoService Suite")
}
