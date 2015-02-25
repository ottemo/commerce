// Copyright 2014 The Ottemo Authors. All rights reserved.

/*
Package db represents database services abstraction layer. It provides a set of interfaces and helpers to interact with
database storage services.

So, when package decides to store something in database it should work with this package instead of trying to interact
with concrete database engine (as Ottemo supposing ability to work with different databases).

Providing Ottemo with a new database engine supposes implementation of "InterfaceDBEngine" with following registration
for db package.
*/
package db

