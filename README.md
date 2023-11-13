THis is sample sqldata application

GETall details
curl -X GET "http://127.0.0.1:8089/data/Getall"

Get Specific country details
curl -X GET "http://127.0.0.1:8089/data/?country="Iran""

Update specific country details
curl -X PUT "http://127.0.0.1:8089/data/SriLanka" -d '{"Runs": 278, "Overs": 50.0}'

Delete specific country details
curl -X DELETE "http://127.0.0.1:8089/data/?country="Pakisthan""

Add country details
curl -X POST "http://127.0.0.1:8089/data/" -d '{"Country":"Iraq", "Runs": 328, "Overs": 30.0}'
