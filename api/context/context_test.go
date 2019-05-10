package context

import (
	"testing"
	"fmt"
	"sync"
)

func ExampleRunInContext() {
	B := func() {
		fmt.Println(GetContext()["test"])
	}

	A :=  func() {
		dict := GetContext()
		if val, ok := dict["x"].(int); ok {
			dict["x"] = val * 2
		}
		B()
	}

	RunInContext(func() { A() }, map[string]interface{} {"test": 1})
}

func TestSimple(t *testing.T) {

	B := func(testValue interface{}) {
		if GetContext()["test"] != testValue {
			t.Fatalf("invalid value")
		}
	}

	A := func(testValue interface{}) {
		fmt.Println(GetContextId())
		if context := GetContext(); context != nil {
			context["test"] = testValue
		} else {
			t.Fatalf("no context found")
		}

		B(testValue)
	}

	MakeContext(func() { A(1) })

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			MakeContext(func() { A(i) })
		}()
	}
	wg.Wait()
}
