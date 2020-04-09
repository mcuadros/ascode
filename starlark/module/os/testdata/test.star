load('os', 'os')
load('assert.star', 'assert')

assert.ne(os.temp_dir(), "")

key = "HELLO"
value = "world"

os.setenv(key, value)
assert.eq(os.getenv(key), value)

content = "hello world\n"

home = os.getenv("HOME")
os.chdir(home)
assert.eq(os.getwd(), home)

temp = os.temp_dir()
os.chdir(temp)

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

# test command
temp = os.temp_dir() + "/example-dir"
os.mkdir_all(temp + "/foo/bar", 0o755)
os.chdir(temp)

assert.eq(os.command("ls -1"), "foo")

# test command dir
assert.eq(os.command("ls -1", dir="foo"), "bar")

# test command shell and env
assert.eq(os.command("echo $FOO", shell=True, env=["FOO=foo"]), "foo")

# test command combined
assert.ne(os.command("not-exists || true", shell=True, combined=True), "")

# test command error
assert.fails(lambda: os.command("not-exists"), "executable file not found")

os.remove_all(temp)
