### Configuration

Set BADSEC_ENDPOINT environment variable to change the API endpoint. As default (if not set) the BADSEC endpoint is localhost:8888.

### How to run it from source

To build from source you will need go1.11 and Docker

The BADSEC server will run on port 8888. You can start it with:

``` docker run --rm -p 8888:8888 adhocteam/noclist ``` 

Build and execute the program with:

``` go build && ./asapp-noclist ``` 

### Run it using Docker

Start BADSEC server

``` docker run --rm -p 8888:8888 adhocteam/noclist ``` 

Build docker container

``` docker build -t asapp-noclip . ```

Run one time docker container and get the results

``` docker run --network="host" -it asapp-noclip ```