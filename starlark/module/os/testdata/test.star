load('os', 'os')
load('assert.star', 'assert')

assert.ne(os.temp_dir(), "")

key = "HELLO"
value = "world"

os.setenv(key, value)
assert.eq(os.getenv(key), value)

content = "hello world\n"

temp = os.temp_dir()
os.chdir(temp)
assert.eq(os.getwd(), temp)

assert.eq(os.write_file("plain.txt", content), None)
assert.eq(os.read_file("plain.txt"), content)

assert.eq(os.write_file("perms.txt", content=content, perms=0o777), None)
assert.eq(os.read_file("perms.txt"), content)

os.mkdir("foo", 0o755)
os.remove("foo")

os.mkdir_all("foo/bar", 0o755)

os.rename("foo", "bar")
os.remove_all("bar")

def deleteNotExistant(): os.remove("foo")
assert.fails(deleteNotExistant, "remove foo: no such file or directory")