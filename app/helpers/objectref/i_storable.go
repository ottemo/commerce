package objectref

import "errors"



// returns current object id
func (it *DBObjectRef) GetId() string {
	return it.id
}



// sets new object id for current object
func (it *DBObjectRef) SetId(newId string) error {
	it.id = newId

	return nil
}


func (it *DBObjectRef) Save() error {
	return errors.New("not implemented")
}

func (it *DBObjectRef) Load(id string) error {
	return errors.New("not implemented")
}

func (it *DBObjectRef) Delete(id string) error {
	return errors.New("not implemented")
}
