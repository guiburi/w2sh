package w2sh

import "text/template"

var tmpl, _ = template.New("webpage").Parse(`
{{ define "message" }}
	{{ range $key, $value := .Cmds }}
		<li><a href="/{{ $key }}">{{ $key }}</a></li>
	{{ end }}
	{{ if gt (len .Children) 0}}
		<ol>
			{{ range .Children }}
				{{ template "message" . }}
			{{ end }}
		</ol>
	{{ end }}
{{ end }}
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
		<ol>{{ template "message" .CmdMap }}</ol>
	</nav>
	<h1>{{.Title}}</h1>
	{{.Form}}
	<br>
	<code style="white-space: pre-line">{{.Output}}</code>
	</body>
</html>
`)
