package main

import "html/template"

var IssueList = template.
	Must(template.New("issueList").
		Parse(`
<h1>Users Table</h1>
	<table>
		<tr style='text-align: left'>
			<th>Name</th>
			<th>Age</th>
			<th>Role</th>
		</tr>
	{{range .}}
		<tr>
			<td><b>{{.Name}}</b></td>
			<td>{{.Age}}</td>
			<td>{{.Role}}</td>
		</tr>
	{{end}}
	</table>
		`))
