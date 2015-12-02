
ENDPOINT = "https://s3.amazonaws.com/ypncks/data.json"

$.getJSON(ENDPOINT, function (data) {
	console.log(data)
})