# yelp-parser

This project is a web go scraper for reviews specifically for Yelp to get reviews for each subject and reviews. This scraper is written in GO to utilize it's robust and easy to use concurrency/parellel framework. 

The parser looks for an object in S3 that has a list of URLS from yelp and will scrape all those URLS in parallel using go channels for concurrent syncs and write into a map interface to write back up to S3 for review analysis.

Dependencies:

* [GoLang](https://golang.org/)
* [Poly GO Utils library](https://github.com/polyglotDataNerd/poly-Go-utils)
* [GoLang SDK for AWS](https://sg.com/yelp/goaws.amazon.com/sdk-for-go/)
* [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)
* [Yelp](https://www.yelp.com/)

