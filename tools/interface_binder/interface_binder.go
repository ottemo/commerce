package interface_binder

import "errors"


var defaultInterfaceBinder *DefaultInterfaceBinder

func init() {
	defaultInterfaceBinder = new(DefaultInterfaceBinder)
	defaultInterfaceBinder.interfaces = map[string]InterfaceCandidate{}
	defaultInterfaceBinder.candidates = map[string][]InterfaceCandidate{}
}

func GetInterfaceBinder() *DefaultInterfaceBinder {
	return defaultInterfaceBinder
}


type ConstructorFunc func() interface{}

type InterfaceCandidate struct {
	CandidateName string
	InterfaceName string
	ConstructorFunc ConstructorFunc
}

type DefaultInterfaceBinder struct {
	interfaces map[string]InterfaceCandidate
	candidates map[string][]InterfaceCandidate
}

func (it *DefaultInterfaceBinder) RegisterCandidate(CandidateName string, InterfaceName string, ConstructorFunc ConstructorFunc, IsDefault bool) error {
	thisCandidate := InterfaceCandidate{ CandidateName: CandidateName, InterfaceName: InterfaceName, ConstructorFunc: ConstructorFunc }
	if _, present := it.candidates[InterfaceName]; !present {

		if cap(it.candidates[InterfaceName]) == 0 {
			it.candidates[InterfaceName] = make([]InterfaceCandidate, 0, 10)
		}

		it.candidates[InterfaceName] = append(it.candidates[InterfaceName], thisCandidate )
	} else {
		return errors.New("candidate with name [" + CandidateName + "] already registered")
	}

	if _, present := it.interfaces[InterfaceName]; IsDefault && !present {
		it.interfaces[InterfaceName] = thisCandidate
	} else {
		return errors.New("default candidate for interface [" + InterfaceName + "] already registered")
	}

	return nil
}

func (it *DefaultInterfaceBinder) GetObject(InterfaceName string) interface{} {
	if candidate, present := it.interfaces[InterfaceName]; present {
		if candidate.ConstructorFunc != nil {
			return candidate.ConstructorFunc()
		} else {
			return nil //, errors.New("constructor function was not defined for candidate [" + candidate.CandidateName + "]")
		}

	} else {
		if len(it.candidates[InterfaceName]) > 0 {
			it.interfaces[InterfaceName] = it.candidates[InterfaceName][0]
			return it.GetObject(InterfaceName)
		}
		return nil //, errors.New("constructor was not found for interface [" + InterfaceName + "]")
	}
}

func (it *DefaultInterfaceBinder) ListInterfaces() []string {
	result := make([]string, 0, len(it.interfaces))
	for key, _ := range it.interfaces {
		result = append(result, key)
	}
	return result
}

func (it *DefaultInterfaceBinder) ListCandidates(InterfaceName string) []string {
	if candidates, present := it.candidates[InterfaceName]; present {
		result := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			result = append(result, candidate.CandidateName)
		}
		return result
	}

	return make([]string, 0, 0)
}
