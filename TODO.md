# TODO

* Unit tests
* Support for dynamic config files based on environment, or pass through config file as os.Args
* Support for other Airbrake notification versions:
  * e.g with subtle differences such as 2.2, 2.1?
* Better handling when invalid requests received?
  * Empty body entirely
  * When certain fields missing?
* Review/double-check and ensure all information is uploaded to Sentry
* Load test - what happens at 100 concurrent connections?
* Find a way to better detail with errors, e.g not panic() in http handlers?
