load("assert.star", "assert")

helm = tf.provider("helm", "1.0.0", "default")
helm.kubernetes.token = "foo"

# require scalar arguments
helm.resource.release()
errors = validate(helm)
assert.eq(len(errors), 2)
assert.eq(errors[0].pos, "testdata/validate.star:7:22")
assert.eq(errors[1].pos, "testdata/validate.star:7:22")

# require list arguments
google = tf.provider("google")
r = google.resource.organization_iam_custom_role(role_id="foo", org_id="bar", title="qux")
r.permissions = ["foo"]
assert.eq(len(validate(google)), 0)

r.permissions.pop()
assert.eq(len(validate(google)), 1)

# require blocks
google = tf.provider("google")
r = google.resource.compute_global_forwarding_rule(target="foo", name="bar")
r.metadata_filters()
assert.eq(len(validate(google)), 2)

errors = validate(google)
for e in errors: print(e.pos, e.msg)