package w2sh

import "text/template"

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<link href="data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII="
          rel="icon" type="image/x-icon"/>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
	<nav>
		<ol>
			{{ range $key, $value := .CommandMap }}
				<ul><a href="/{{ $key }}">{{ $key }}</a></ul>
			{{ end }}	
		</ol>
	</nav>
	<h1>{{.Title}}</h1>
	{{.Form}}
	<br>
	<code style="white-space: pre-line">{{.Output}}</code>
	</body>
</html>`

var tmpl, _ = template.New("webpage").Parse(tpl)
