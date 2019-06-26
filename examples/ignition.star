ignition = provider("ignition")
print("provider  --->", dir(ignition))

user = ignition.user("test")
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True



disk = ignition.disk("foo")
disk.device = "/dev/sda"

root = disk.partition("root")
root.start = 2048
root.size = 4 * 1024 * 1024 

home = disk.partition("home")
home.start = root.size + root.start 
home.size = 4 * 1024 * 1024 

print(home.start)