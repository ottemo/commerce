package otto
//
//import (
//	"strings"
//	"github.com/chzyer/readline"
//	"github.com/robertkrimen/otto/repl"
//)
//
//func StartCli2() {
//	repl.Run(baseVM)
//}
//
//func StartCli() error {
//	vm := baseVM.Copy()
//	readLine, err := readline.NewEx(&readline.Config{Prompt: "otto> "})
//
//	if err != nil {
//		panic(err)
//	}
//
//	defer readLine.Close()
//
//	script := ""
//	for {
//		line, err := readLine.Readline()
//		if err != nil {
//			return err
//		}
//		line = strings.TrimSpace(line)
//		script += line
//
//
//		if line == "" {
//			continue;
//		}
//
//		readLine.SaveHistory(script)
//
//		if !strings.HasSuffix(line, ";") &&
//			strings.Count(script, "{") == strings.Count(script, "}") &&
//			strings.Count(script, "(") == strings.Count(script, ")") &&
//			strings.Count(script, "[") == strings.Count(script, "]") {
//
//			result, err := vm.Eval(script)
//			if err != nil {
//				readLine.Terminal.Print(err.Error())
//				readLine.Terminal.Print("\n")
//			}
//
//			value, err := result.ToString()
//			if err != nil {
//				readLine.Terminal.Print(err.Error())
//				readLine.Terminal.Print("\n")
//			}
//
//			readLine.Terminal.Print(value)
//			readLine.Terminal.Print("\n")
//
//			script = ""
//			readLine.SetPrompt("otto> ")
//			continue
//		} else {
//			readLine.SetPrompt("    > ")
//		}
//
//		script += "\n"
//	}
//
//	return nil
//}
//
//
//func registerValue(name string, value interface{}) {
//	baseVM.Set(name, value)
//}
