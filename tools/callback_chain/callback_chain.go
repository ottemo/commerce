package callback_chain

import (
	"errors"
)


// Types and Variables

type T_CallbackValidator func(interface{}) bool
type T_CallbackChainExecutor func( []T_CallbackFunc )

type T_CallbackFunc struct {
	   Alias string
	Function interface{}
}

type T_CallbackChainItem struct {
	         Name string
	    Validator T_CallbackValidator
	     Executor T_CallbackChainExecutor
	CallbackChain []T_CallbackFunc
}

var callbackChains = map[string]T_CallbackChainItem {}


// Stub callback functions

func DefaultVoidFuncWithNoParamsValidator(Input interface{}) bool {
	_, ok := Input.(func())
	return ok
}

func DefaultVoidFuncWithNoParamsExecutor( function interface{} ) {
	function.(func())()
}

func DefaultChainVoidFuncWithNoParamsExecutor( chain []T_CallbackFunc ) {
	for _, callbackFunc := range chain {
		callbackFunc.Function.(func())()
	}
}


// Functionality

func RegisterCallbackChain(Name string, Validator T_CallbackValidator, Executor T_CallbackChainExecutor) error {
	if _, present := callbackChains[Name]; present {
		return errors.New("callback chain '" + Name + "' already registered")
	} else {
		callbackChains[Name] = T_CallbackChainItem{Name: Name, Validator: Validator, Executor: Executor, CallbackChain: make([]T_CallbackFunc, 0) }
	}

	return nil
}

func RegisterCallbackForChain(ChainName string, TargetAlias string, BeforeTarget bool, Alias string, CallbackFunc interface{}) error {
	// checking chain name is valid
	if callbackChainItem, present := callbackChains[ChainName]; present {
		// checking callback function is valid
		if callbackChains[ChainName].Validator(CallbackFunc) {

			callbackChain := callbackChainItem.CallbackChain

			newOneCallback := T_CallbackFunc { Alias: Alias, Function: CallbackFunc }
			newCallbackChain := make([]T_CallbackFunc, len(callbackChain)+1)

			newIndex := 0
			found := false
			for _, callbackFunc := range callbackChainItem.CallbackChain {
				// we are not allowing to use same Alias within chain
				if callbackFunc.Alias == Alias {
					return errors.New("alias '" + Alias + "' already exists in chain '" + ChainName + "'")
				}

				// TargetAlias was found
				if callbackFunc.Alias == TargetAlias {
					if BeforeTarget {
						newCallbackChain[newIndex] = newOneCallback
						newCallbackChain[newIndex+1] = callbackFunc
						newIndex = newIndex + 2
					} else {
						newCallbackChain[newIndex] = callbackFunc
						newCallbackChain[newIndex+1] = newOneCallback
						newIndex = newIndex + 2
					}
					found = true
					continue
				}
				newCallbackChain[newIndex] = callbackFunc
				newIndex = newIndex + 1
			}

			// Target Alias was not found - appending to the end
			if !found {
				callbackChain[len(callbackChain)-1] = newOneCallback
			}

		} else {
			return errors.New("not valid callback function")
		}
	} else {
		return errors.New("can not find callback chain '" + ChainName + "'")
	}

	return nil
}

func RegisterCallbackInChainBefore(ChainName string, Before string, Alias string, CallbackFunc interface{}) error {
	return RegisterCallbackForChain(ChainName, Before, true, Alias, CallbackFunc)
}
func RegisterCallbackInChainAfter(ChainName string, After string, Alias string, CallbackFunc interface{}) error {
	return RegisterCallbackForChain(ChainName, After, false, Alias, CallbackFunc)
}
func RegisterCallbackInChainTail(ChainName string, Alias string, CallbackFunc interface{}) error {
	return RegisterCallbackForChain(ChainName, "", false, Alias, CallbackFunc)
}

func ExecuteCallbackChain(ChainName string) {
	if callbackItem, present := callbackChains[ChainName]; present {
		callbackItem.Executor( callbackItem.CallbackChain )
	}
}
