load('os', 'os')
load('assert.star', 'assert')

key = "HELLO"
value = "world"

os.setenv(key, value)
assert.eq(os.getenv(key), value)

content = "hello world\n"

assert.eq(os.write_file("/tmp/plain.txt", content), None)
assert.eq(os.read_file("/tmp/plain.txt"), content)

assert.eq(os.write_file("/tmp/perms.txt", content=content, perms=0o777), None)
assert.eq(os.read_file("/tmp/perms.txt"), content)

os.mkdir("foo", 0o755)
os.remove("foo")

assert.ne(os.getwd(), "")

home = os.getenv("HOME")
os.chdir(home)
assert.eq(os.getwd(), home)

os.mkdir_all("foo/bar", 0o755)

os.rename("foo", "bar")
os.remove_all("bar")

def deleteNotExistant(): os.remove("foo")
assert.fails(deleteNotExistant, "remove foo: no such file or directory")