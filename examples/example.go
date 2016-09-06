package examples

func GetWeatherData() []byte {
	// Data returned from
	// https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%20from%20weather.forecast%20where%20woeid%20%3D%202487889&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys
	return []byte(`
		{
		    "query": {
			"count": 1,
			"created": "2016-09-06T17:56:20Z",
			"lang": "en-GB",
			"results": {
			    "channel": {
				"item": {
				    "condition": {
					"code": "34",
					"date": "Tue, 06 Sep 2016 10:00 AM PDT",
					"temp": "72",
					"text": "Mostly Sunny"
				    }
				}
			    }
			}
		    }
		}
	`)
}
