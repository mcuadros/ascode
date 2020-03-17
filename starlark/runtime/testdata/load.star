# test relative loading
load("includes/foo.star", "foo")

# evaluate and correct base_path
mod = evaluate("includes/foo.star")
print(mod.foo)

# module constructor
module("foo")