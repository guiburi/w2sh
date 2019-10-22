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
	<h1>{{.Title}}</h1>
	{{.Form}}
	<br>
	<code style="white-space: pre-line">{{.Output}}</code>
	</body>
</html>`

var tmplte, err = template.New("webpage").Parse(tpl)
