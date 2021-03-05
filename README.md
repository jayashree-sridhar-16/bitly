# Cloud Project

## Youtube link
https://youtu.be/yQsDZMwYc30

## Extra credit options attempted:
- Hosting a web app on Heroku 
  - A simple UI using ReactJS to create new short links, view previously created links and view the statistics of how many times created links have been accessed
  - Create view configured with Control Panel API created using AWS API gateway for control panel load balancer
  - Trends tab configured with Link Redirect API created using AWS API gateway for link redirect load balancer  
  - Hosted at 
  ```
  https://cmpe281-bitly.herokuapp.com/ 
  ```
  > ( Needs AWS backend instances to be up to work)
  ```
  Reference for creating the app:
  https://www.kirupa.com/react/creating_single_page_app_react_using_react_router.htm
  
  Deploy on heroku:
  https://www.geeksforgeeks.org/how-to-deploy-react-app-to-heroku/
  
  ```
  
- AutoDiscover nosql nodes from autoscale group
  - Code at branch **autodiscover**   
  https://github.com/nguyensjsu/cmpe281-jayashree-sridhar-16/blob/autodiscover/nosql/api/src/main/java/api/NodeDiscovery.java  
  > Code added to separate branch from master to prevent conflict while build nosql without autodiscovery and hosting as separate instances for given deployment model.
  - Used aws sdk java package to get the instances running as part of autoscale group names 'nosql'
  - Get equivalent ec2 instance metadata from instance ids in autoscale group
  - Check status that instances are running using status code
  - If running; check against the nodes map maintained by AdminServer if instance already registered
  - If not, loop through all registered nodes and perform a POST against those nodes to register the newly discovered node
  ```
   Packages used 
      - https://docs.aws.amazon.com/AWSJavaSDK/latest/javadoc/com/amazonaws/services/ec2/model/package-summary.html
      - https://docs.aws.amazon.com/AWSJavaSDK/latest/javadoc/com/amazonaws/services/autoscaling/model/package-summary.html
   ```

## Project Journal
##### Nov 8: 
- Created a Go API for server listening on port 3000 and having routes '/ping' and '/' 
- Added method for testing shortening a url using md5 hashing in golang

##### Nov 10:
- Created database and table in mysql
- Added mysql connection with go app using go-sql-driver package

##### Nov 12 - 15
- Tested redirection api using gorilla/mux

##### Nov 19
- Created separate servers for control panel and link redirect services
- Link redirect server redirecting by reading from mysql backened

##### Nov 20
- Added Trend server, link redirect server and control panel server
- Rabbitmq queues for create from control panel to trend server, create the document of new url in local nosql db cluster

##### Nov 23 - 28
- Tried running NoSQL containers in 2 AWS EC2 instances inside docker
- Changed code in AdminServer.java for ping check of nodes using Socket instead of InetAddress since it was failing for some reason in AWS EC2
- Adding front end web app page structure for Bitly using ReactJS
- Added web page for creating short links and viewing statistics

#### Nov 30
- Hosted all 5 nodes of NoSQL image, in 5 different docker host AWS EC2 instances
- Added a network load balancer in front of the 5 instances
- Tests registering and creating documents against loadbalancer api

#### Dec 1 - 4
- Built docker images for control panel server, trend server and link redirect server
- Pushed the docker images to jay16/controlpanel:aws, jay16/linkredirect:aws and jay16/trendserver:aws 
- Hosted the mentioned images in 3 different docker host aws ec2 instances
- Hosted Rabbitmq in the same instance as trend server
- Tested the basic functionality of the servers communicating with each other without load balancers and autoscaling

#### Dec 5 - 8
- Added network load balancer and autoscale groups for control panel servers and link redirect servers
- Create AWS API Gateway for control panel and link redirect load balancers
- Configured Heroku app to talk with API gateways instead of localhost
- Tested full setup on AWS setup
- Started working on NoSQL autodiscovery using aws sdk java package
- Able to read instances of nosql autoscale group and its intsances; register all instances to each other
- Sync not working with autodiscover, code changes needed

## Implementation details:
### Deployment Diagram:
<img src="https://github.com/nguyensjsu/cmpe281-jayashree-sridhar-16/blob/master/Design/DeploymentDiagram.png" />

#### Deployment steps:
##### Nosql cluster on AWS:
	- Build ap style nosql image and push to docker hub
	- Launch a container from docker-host ami in private subnet and clean up running instances
	- Pull nosql image and run it
	- Create nosql ami and launch 4 more instances of nosql to get 5 node NoSQL cluster in private subnets
	- Create an internal network load balancer against 5 NoSQL nodes
	- SG ports: 80, 22, 8888, 9090
##### Trend server and RabbitMq
	- Launch another docker-host ami in private subnet and clean up running instances
	- Pull RabbitMq docker image and run container based on this image
	- Change rabbitmq host name to instance private IP in trendserver go file
	- Build trendserver go app locally, build docker image and push it to docker hub
	- Pull the image in the ec2 instance and run it with port mapped to 80
	- SG ports: 80, 22, 5672, 8080, 4369
##### Link redirect server
	- Configure to trendserver private IP address 
	- Build link redirect server go file locally, build and push docker image to docker hub
	- Launch another docker-host ami in private subnet and clean up running instances
	- Pull and run linkredirect docker image and port map it to port 80
	- Create ami from this instance, and use it to create launch configuration
	- Create autoscale group from this launch configuration to scale between 1 and 3 instances
	- Add load balancer and target groups for the autoscale group
	- SG ports: 80, 22
##### Control panel server
	- Repeat the same steps for control panel as link redirect server
	- SG ports: 80, 22
##### API Gateway
	- Create REST APIs for control panel load balancer and link redirect load balancer
##### Heroku
	- Configure ReactJS app with created AWS APIs 


### JSON data maintained in NoSQL Key/Value DB
```
type bitly struct {
	Original_url  string
	Short_url  string
	Redirect_url string
	Access_count int
}
```

### Control Panel:
- Runs on port 3000
- Used crypto/md5 golang library for creating encoded or hash url strings (https://golang.org/pkg/crypto/md5/)
- Connects with MySQL database on port 3306 to the database *bitly*
- Connects to RabbitMq server and publishes messages on *createLink* message queue
- APIs:
```
Path: GET /ping
Response HTTP 200 OK
{
    "Test": "Control panel active!"
}

Path: POST /links/create
Request Body:
{
    "Original_url": "https://hub.docker.com/repository/docker/jay16/nosql"
}

Response 200 OK
{
    "Original_url": "https://golang.org/pkg/bytes/",
    "Short_url": "c74c4d36",
    "Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/c74c4d36",
    "Access_count": 0
}
```

### Link Redirect Server:
- Runs on port 3001
- Connects with trend server
- Uses http Redirect for redirecting short url to original url
- Gets original url from trend server
- Gets all links present in trend server and ranks them with links having access count on top
> Used https://golang.org/pkg/sort/#SliceStable for sorting based on access count
```
Reference for redirecting the URL:   
https://restbird.org/docs/mock-examples-golang.html#get-request-header
```
- APIs:
```
Path GET /links
Response 200 OK
[
    {
        "Original_url": "https://www.mercurynews.com/2020/12/07/coronavirus-california-is-scrambling-as-many-regions-begin-the-biggest-shutdown-since-spring/",
        "Short_url": "73b1f242",
        "Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/73b1f242",
        "Access_count": 15
    },
    {
        "Original_url": "https://towardsdatascience.com/gui-fying-the-machine-learning-workflow-towards-rapid-discovery-of-viable-pipelines-cab2552c909f",
        "Short_url": "931d24a4",
        "Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/931d24a4",
        "Access_count": 6
    }    
]

Path GET /redirect/931d24a4
Response 302 Found
```

### Trend Server:
- Runs on port 3002
- Connects to rabbitmq server to consume messages from *createLink* queue
- Coneects to NoSQL db cluster through loadbalancer for creating documents or reading documents
- Documents are created with created hash as key value
- Updates access count each time call for specific document is received and updates the NoSQL db
- APIs:
```
Path GET /links
Response 200 OK
[
    {
        "Original_url": "https://www.mercurynews.com/2020/12/07/coronavirus-california-is-scrambling-as-many-regions-begin-the-biggest-shutdown-since-spring/",
        "Short_url": "73b1f242",
        "Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/73b1f242",
        "Access_count": 15
    },
    {
        "Original_url": "https://towardsdatascience.com/gui-fying-the-machine-learning-workflow-towards-rapid-discovery-of-viable-pipelines-cab2552c909f",
        "Short_url": "931d24a4",
        "Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/931d24a4",
        "Access_count": 6
    }    
]

Path GET /links/931d24a4
Response 200 OK
{
	"Original_url": "https://towardsdatascience.com/gui-fying-the-machine-learning-workflow-towards-rapid-discovery-of-viable-pipelines-cab2552c909f",
	"Short_url": "931d24a4",
	"Redirect_url": "https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/redirect/931d24a4",
	"Access_count": 6
}

```

## Other References:
- Original NoSQL code for ping check using InetAddress.isReachable did not work when NoSQL nodes were hosted in AWS EC2 instances. Used socket for ping check instead and used the code from https://crunchify.com/how-to-implement-your-own-inetaddress-isreachablestring-address-int-port-int-timeout-method-in-java/  as below
```
/*
	 * Overriding default InetAddress.isReachable() method to add 2 more arguments port and timeout value
	 * 
	 * Address: www.google.com 
	 * port: 80 or 443 
	 * timeout: 2000 (in milliseconds)
	 */
	private static boolean crunchifyAddressReachable(String address, int port, int timeout) {
		try {
 
			try (Socket crunchifySocket = new Socket()) {
				// Connects this socket to the server with a specified timeout value.
				crunchifySocket.connect(new InetSocketAddress(address, port), timeout);
			}
			// Return true if connection successful
			return true;
		} catch (IOException exception) {
			exception.printStackTrace();
 
			// Return false if connection fails
			return false;
		}
	}
```
- RabbitMq tutorials  
https://github.com/rabbitmq/rabbitmq-tutorials/tree/master/go
- For basic idea how to implement Bitly  
https://itnext.io/designing-the-shortening-url-system-like-bit-ly-loading-6-billion-clicks-a-month-78b3e48eee8c
- Preflight requests in ReactJS axios issue  
https://weezy.dev/avoiding-preflight-requests-with-axios/
- REST APIs in golang  
https://tutorialedge.net/golang/creating-restful-api-with-golang/
