package interface_binder_test

import ( "testing"
		 "reflect"
		 "github.com/ottemo/platform/tools/interface_binder"
)

type I_A interface {
	Init(int, int)
	Multiply() int

	InitB(int, int)
	MultiplyB() int
}

type A struct {
	x int
	y int

	B
}

type B struct {
	xx int
	yy int
}

type C struct {
	A
}

func (it *C) Multiply() int {
	return it.x
}

func (it *B) InitB(xx int, yy int) {
	it.xx = xx
	it.yy = yy
}

func (it *B) MultiplyB() int {
	return it.xx * it.yy
}

func (it *A) Init(x int, y int) {
	it.x = x
	it.y = y
}

func (it *A) Multiply() int {
	return it.x * it.y
}



func AConstructor() interface{} {
	return new(A)
}

func BConstructor() interface{} {
	return new(B)
}

func TestInterfaceBinder0(t *testing.T) {
	x := new(C)
	x.Init(10, 2)
	t.Log( x.A.Multiply() )
}

func TestInterfaceBinder1(t *testing.T) {
	ib := interface_binder.GetInterfaceBinder()

	ib.RegisterCandidate("Default", "I_A", AConstructor, false)
	ib.RegisterCandidate("B", "I_A", BConstructor, true)

	if obj, ok := ib.GetObject("I_A").(I_A); ok {
		obj.Init(10, 20)
		result := obj.Multiply()

		obj.InitB(2, 2)
		result2 := obj.MultiplyB()

		if result == 10*20 && result2 == 2*2 {
			t.Log("10 * 20 = ", result)
			t.Log("2 * 2 = ", result2)
		}else{
			t.Error("wrong calculation - something wrong")
		}
	}else{
		t.Error("can't get interface candidate", reflect.TypeOf( ib.GetObject("I_A") ))
	}

}
