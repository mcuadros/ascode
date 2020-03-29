# test relative loading
load("includes/foo.star", "foo")

# evaluate and correct base_path
mod = evaluate("includes/foo.star")
print(mod.foo)

# module constructor
module("foo")

# test defined modules
load("encoding/json", "json")
load("encoding/base64", "base64")
load("encoding/csv", "csv")
load("encoding/yaml", "yaml")
load("math", "math")
load("re", "re")
load("time", "time")
load("http", "http")