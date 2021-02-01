tmpl = tools.newTmpl()
tmpl.parse("t1", [[
t1 echo {{.}}
]])
tmpl.parse("t2", [[
t2 echo {{.name}}
]])

s1 = tmpl.execute("t1", 'panda')
s2 = tmpl.execute("t2", {name="pandacc",})
print(s1)
print(s2)