# poly-yelp

This project is a web go scraper for reviews specifically for Yelp to get reviews for each subject and reviews. The scraper is written in GO to utilize it's robust and easy to use concurrency/parellel framework. 

The parser looks for an object in S3 that has a list of URLS from yelp and will scrape all those URLS in parallel using go channels for concurrent syncs and write into a map interface to write back up to S3 for review analysis.

**Dependencies**

* [GoLang](https://golang.org/)
* [Poly GO Utils library](https://github.com/polyglotDataNerd/poly-Go-utils)
* [GoLang SDK for AWS](https://sg.com/yelp/goaws.amazon.com/sdk-for-go/)
* [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)
* [Yelp](https://www.yelp.com/)

Application Arguments
-

| Argument        | Sample           | Required  |
| ------------- |:-------------:| -----:|
| s3 input | Optional input param, is used when reading a file in s3 rather than an URL  | NO  |
| s3 output | landing s3 directory of data i.e. s3://target | YES  |
| key output | target key of output bucket i.e. s3://target **/etc** | YES  |
| load type | **flat**: s3 object with all url paths, **url**: looks for one URL only | YES  |
| source urls | **load type = flat**: recursion on s3 dir, **load type = url**: passes one URL only | YES  |

`Noteable Mention On Args`
 1. **load type** is dependant on how the user wants to run the application
    
    * _flat_: takes multiple objects in an s3 directory of URL addresses puts into many go routines and runs a parallel process in a container that pings and scrapes all URLS into a single s3 object. 
    * _url_: passes a single URL for the container to run, each container will output an s3 object. This pattern uses many containers that can run parallel with independent resources; rather than use one machine it uses a cluster of many using the same docker image with a different URL in the args. 

2. **source urls** is dependant on how the user wants to run the application
    
    * _flat_: an s3 path **s3://target/etc** populated with URLS
    * _url_: a single URL https://www.yelp.com/biz/mountain-cafe-los-angeles-4
