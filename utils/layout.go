package utils

const layout = `<html>
<head>
	<style>
	th, td {
	  border: 1px solid #ddd;
	  text-align: left;
	  padding: 8px;
	}
	th {
		padding: 8px;
		background-color: #4CAF50;
	}
	tr:nth-child(even){background-color: #f2f2f2;}
	tr:hover {background-color:#f5f5f5;}
	table {
		border: 1px solid #ddd;
		width:70%; 
margin-left:15%; 
margin-right:15%;
}
img {
display: block;
margin-left: 35%;
margin-right: 35%;
}
	</style>
	</head>
	<body>
		<h1>{{.PageTitle}}</h1>
		<div>
			<img src="results_pie_chart.png" alt="Test Results" style="width:30%">
		</div>
<div>
<table>
	<tr>
	  <th>File Name</th>
	  <th>Passed</th> 
	  <th>Failed</th>
	  <th>Ignored</th>
	  <th>Total</th>
	  <th>Status</th>
	</tr>
	{{range .Tests}}
	<tr>
			<td>{{.FileName}}</td>
			<td>{{.Passed}}</td>
			<td>{{.Failed}}</td>
			<td>{{.Ignored}}</td>
			<td>{{.Total}}</td>
			<td>{{.Status}}</td>
	</tr>
	{{end}}
  </table>
</div>
	</body>
</html>`
