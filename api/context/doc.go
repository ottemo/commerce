/*
Package context implements call-stack context storage, it allows to store and pass the objects/values for nested
function calls regard to application call-stack. Use it whenever you need to pass value(s) to nested function calls
where it is impossible or inconvenient to pass values through function arguments.

Example usage:

Imagine writing a method to an HTTP application endpoint that is quite complicated. That endpoint is mapped to call
function "A()", which then performs calls to functions B, C, D. Then to add complexity, they have nested-calls to
E, F, G, H as well, like:

  A(request, response) -> B(request) -> E() -> F() -> H()
                       -> C() -> G()
                       -> D()

There exists a use case where a warning could happen within function "E()" and this information is significant up the
call stack (i.e. you need to pass that warning to function "A()"). You can update function B and E with the reference
argument or update the E function to pass warning with the returning value - but this will require appropriate
modifications in "parent" functions "A()" and "B()". Alternatively you can use a global package variable and put the
message there, but a global package variable would stay the same for different "go-routines" and the function "A()"
would not have a chance to distinguish the request relationship to the data.

Consider that "H()" function (or any more nested function) could exemplify the similar case or for an instance the "H()"
function is located in an external package and works with the interface but it requires HTTP request parameters to perform
the appropriate work, additionally consider the fact these parameters are not static.

If this has not convinced you, consider the following case: you are developing a database driver in a separate package, the
other packages will work with your driver through interfaces. The database supports the transactions. There are at
least 2 packages (shipping and payment) which are using the DB operations during the checkout process (another package).
Shipping and payment packages should save the generated information into database only in the case when the checkout is
successful and the order is created, otherwise the operations should be rolled back. All packages are being developed
independently by different organizations.

  [checkout]: receipt() -> [db] startTransaction()
                        -> [shipping] createShipping(...) -> [db] insertData(...)
                        -> [payment] createPayment(...) -> [db] insertData(...)
                        -> [checkout] placeOrder(...) -> error()
                        -> [db] rollbackTransaction()

The same database driver is used for different operations and the several DB operations happen simultaneously in
different go routines (cms operation, another checkout, etc.). How the database driver could distinguish the operations
it should rollback on "rollbackTransaction()" call? The transaction ID should be passed to each possible package and
later provided for each possible functions working with DB? Should the package be aware of the fact the transaction
occurs or should the pakcage leave this level of understanding of the call stack to the database driver and checkout
packages?

We suggest you address the architecture nightmare using this context package.

How does this work?

There is a call stack you have access to in GO applications via a call to runtime.CallersFrames(). This method returns the
frames created by Golong for each function call. It follows:

      {PC}:{Entry}  - {File}:{Line}                - {Function}
  -------------------------------------------------------------------------
  17392491:17392432 - /Users/user/go/sample.go:421 - main.H
  17392415:17392384 - /Users/user/go/sample.go:417 - main.F
  17392367:17392336 - /Users/user/go/sample.go:413 - main.E
  17392319:17392288 - /Users/user/go/sample.go:409 - main.B
  17392260:17392208 - /Users/user/go/sample.go:405 - main.A
  17392768:17392688 - /Users/user/go/sample.go:425 - main.main
  16949355:16948832 - /usr/local/go/src/runtime/proc.go:200 - runtime.main

Each stack frame will have the "Entry" field which is a pointer to the function entity. That pointer is unique and persistent
for each function in the application.

That is how the above will change when using the suggested the context package.

  17392491:17392432 - /Users/user/go/sample.go:421 - main.H
  17392415:17392384 - /Users/user/go/sample.go:417 - main.F
  17392367:17392336 - /Users/user/go/sample.go:413 - main.E
  17392319:17392288 - /Users/user/go/sample.go:409 - main.B
  17392260:17392208 - /Users/user/go/sample.go:405 - main.A
  17409360:17409280 - /Users/user/go/sample.go:425 - main.main.func1          <- c.RunInContext(func() { A() }, nil)
  17401636:17401600 - /Users/user/go/context.go:150 - context.glob..func139   <- proxiesDict[17401636] = 139
  17389962:17389840 - /Users/user/go/context.go:302 - context.proxyLoop       <- proxiesDict[17389840] = 256
  17389915:17389840 - /Users/user/go/context.go:300 - context.proxyLoop       <- proxiesDict[17389840] = 256
  17392034:17391696 - /Users/user/go/context.go:386 - context.RunInContext    <- proxiesStart
  17392741:17392688 - /Users/user/go/sample.go:425 - main.main
  16949355:16948832 - /usr/local/go/src/runtime/proc.go:200 - runtime.main

You can see the call stack, which allows us to bind unique number/key and use the kv pair like a global variable. In the
application we wrap the call to function "A()" with context.RunInContext(func() { A() }, nil) which looks for a
free value in the global "contexts" variable - this happens using a regular loop from 0 to uint max value, and once the
free key is found in the "contexts" variable it inserts the new "map[string]interface{}" instance which is the context.
Note, RunInContext has a defer function which will release the key once the given function is executed. In the above example
this key is derived in ths manner: "256+256+139=651".

The context package has 255 similar function and knows the corresponding pointers, each pointer is mapped to a
sequenced number from 1 to 255, if the index is bigger than 255 it calls the recursive function "proxyLoop" reducing
the index by 256 each time until there is a remainder less than 256 which can be mapped onto existing functions, (i.e.
less then proxyBase which equals length of "proxies" variable, currently it is defined as 255 proxy functions).

Example:
    import "context"
    import "fmt"

    func A() {
            fmt.Printf("func A: x=%d\n", GetContext()["x"])
    }

    func B() {
          dict := GetContext()
          fmt.Printf("func B: x=%d\n", dict["x"])

            if val, ok := dict["x"].(int); ok {
                dict["x"] = val * 2
            }

            A()
    }

    func C() {
          dict := GetContext()
          fmt.Printf("func C: x=%d\n", dict["x"])
        dict["x"] = 5

            B()
    }

    func main() {
            RunInContext(C, map[string]interface{} {"x": -1})
    }


Result:
    func C: x=-1
    func B: x=5
    func A: x=10

*/
package context
