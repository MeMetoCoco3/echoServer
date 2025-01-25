package main

import "html/template"

var GetAllUsers = template.
	Must(template.New("issueList").
		Parse(
			`<!DOCTYPE html>
<html>
<head>
	<title>User Management</title>
		<link rel="stylesheet" href="/static/style.css">
	</head>
<body>
	<h1>Users Table</h1>
	<table>
		<tr>
			<th>ID</th>
			<th>Name</th>
			<th>Age</th>
			<th>Role</th>
			<th>Actions</th>
		</tr>
		{{range .}}
		<tr>
			<td class="uuid-td">{{.ID}}</td>
			<td><b>{{.Name}}</b></td>
			<td>{{.Age}}</td>
			<td>{{.Role}}</td>
			<td>
		<form action="/delete/{{.ID}}" method="POST" class="delete-form">
					<button type="submit" class="delete-btn">Delete</button>
				</form>
			</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`))
