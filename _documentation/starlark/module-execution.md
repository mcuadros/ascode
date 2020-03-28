---
title: 'Module execution'
weight: 8
---

Each Starlark file defines a _module_, which is a mapping from the
names of global variables to their values.
When a Starlark file is executed, whether directly by the application
or indirectly through a `load` statement, a new Starlark thread is
created, and this thread executes all the top-level statements in the
file.
Because if-statements and for-loops cannot appear outside of a function,
control flows from top to bottom.

If execution reaches the end of the file, module initialization is
successful.
At that point, the value of each of the module's global variables is
frozen, rendering subsequent mutation impossible.
The module is then ready for use by another Starlark thread, such as
one executing a load statement.
Such threads may access values or call functions defined in the loaded
module.

A Starlark thread may carry state on behalf of the application into
which it is embedded, and application-defined functions may behave
differently depending on this thread state.
Because module initialization always occurs in a new thread, thread
state is never carried from a higher-level module into a lower-level
one.
The initialization behavior of a module is thus independent of
whichever module triggered its initialization.

If a Starlark thread encounters an error, execution stops and the error
is reported to the application, along with a backtrace showing the
stack of active function calls at the time of the error.
If an error occurs during initialization of a Starlark module, any
active `load` statements waiting for initialization of the module also
fail.

Starlark provides no mechanism by which errors can be handled within
the language.

