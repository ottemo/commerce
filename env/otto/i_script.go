package otto

import (
	"github.com/robertkrimen/otto/repl"
	"github.com/chzyer/readline"
	"strings"
)

func (it *Script) Execute(code string) (interface{}, error) {
	return it.vm.Eval(code)
}

func (it *Script) Get(name string) (interface{}, error) {
	return it.vm.Get(name)
}

func (it *Script) Set(name string, value interface{}) error {
	return it.vm.Set(name, value)
}

func (it *Script) Interact2() error {
	return repl.Run(it.vm)
}

func (it *Script) Interact() error {
	readLine, err := readline.NewEx(&readline.Config{Prompt: "otto> "})

	if err != nil {
		panic(err)
	}

	defer readLine.Close()

	script := ""
	for {
		line, err := readLine.Readline()
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		script += line

		if line == "" {
			continue;
		}

		readLine.SaveHistory(script)

		if !strings.HasSuffix(line, ";") &&
			strings.Count(script, "{") == strings.Count(script, "}") &&
			strings.Count(script, "(") == strings.Count(script, ")") &&
			strings.Count(script, "[") == strings.Count(script, "]") {

			result, err := it.vm.Eval(script)
			if err != nil {
				readLine.Terminal.Print(err.Error())
				readLine.Terminal.Print("\n")
			}

			value, err := result.ToString()
			if err != nil {
				readLine.Terminal.Print(err.Error())
				readLine.Terminal.Print("\n")
			}

			readLine.Terminal.Print(value)
			readLine.Terminal.Print("\n")

			script = ""
			readLine.SetPrompt("otto> ")
			continue
		} else {
			readLine.SetPrompt("    > ")
		}

		script += "\n"
	}

	return nil
}
