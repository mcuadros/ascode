# Evaluate execute the given file, with the given context, in this case `foo`
values = evaluate("evaluable.star", foo="foo")

# The context it's a module, in this case contains the key `bar`
print("Print from main: '%s'" % values.bar)

# Output:
# Print from evaluable.star: 'foo'
# Print from main: 'bar'
