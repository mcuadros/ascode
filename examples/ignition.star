ignition = provider("ignition")

user = ignition.data.user("test")
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True

print(user.__dict__)

disk = ignition.data.disk("foo")
disk.device = "/dev/sda"

root = disk.partition("root")
root.start = 2048
root.size = 4 * 1024 * 1024 

home = disk.partition("home")
home.start = root.size + root.start 
home.size = 4 * 1024 * 1024 

print("parition count: ", len(disk.partition))
print(disk.__dict__)