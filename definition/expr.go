package definition

type Expr struct {
	Lang       string
	Expression string
}

func Simple(expression string) Expr {
	return Expr{Lang: "simple", Expression: expression}
}

func JsonQuery(query string) Expr {
	return Expr{Lang: "jsonquery", Expression: query}
}

// Template returns a definition of the golang template
func Template(textTemplate string) Expr {
	return Expr{Lang: "template", Expression: textTemplate}
}

// HTML returns a definition of the golang html template
func HTML(htmlTemplate string) Expr {
	return Expr{Lang: "html", Expression: htmlTemplate}
}

// LUA return a definition of the LUA script
func LUA(luaScript string) Expr {
	return Expr{Lang: "lua", Expression: luaScript}
}

func JMESPath(jmesPath string) Expr {
	return Expr{Lang: "jmespath", Expression: jmesPath}
}
