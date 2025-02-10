# Fetch Backend API - Receipt Processor
## Coded by Jesse Jones

# Building Docker Image
## To create the Docker image for this program:

 0. Be on a unix-like environment. If you're on Windows, use the Dockerfile 
    and Go files to build it using some different method.
 1. Change directory into `webserver` directory
 2. Run the bash script `buildImage.sh`, which will run the commands necessary 
    to generate a Docker image based on the Dockerfile and the three Go files provided.
    This will require sudo privelages to run as is. If the user running this has permissions for Docker,
    run the command: `docker build -t server .`, which does the same thing as the bash script but doesn't
    use sudo at all. 
 3. That's all! The Docker image is now built!

# Running the Docker Image and Starting Server
## To run the server:
 
 0. Be on a unix-like environment. If on something like Windows, 
    figure out how to start up the container some other way, maybe via the GUI, 
    or some windows version of the following steps.
 1. Run the bash script `runServer.sh`, which starts the server up, with it lisenting at port 8000.
 2. As with building, if you don't want to use sudo and have privelages to use Docker, 
    simply type the command: `docker run -p 8000:8000 server`, which starts up the server the same way
    as the script does just without sudo.
 3. Once the container is running, whatever terminal started it will see a message printed saying:
    > Listening on port 8000!
    
    Indicating that the server is ready to begin processing queries.
 4. If port 8000 doesn't work, simply change the Dockerfile and code to use a different port number.

# Testing the Server
Beyond the automatic testing that will be done, a few example jsons are provided with
a couple bash scripts that can insert them into the database via a curl request or
query the server with a given ID. 

The script `insertReceiptToDb` with an argument of a JSON file, is what's used
to more easily insert a JSON receipt into the database.

The script `queryPointsAtID.sh` takes an ID that can be found from the output
of the server, and produces the point count associated with it.

# Stopping the Server
To stop the server, just stop the container from running via killing it directly, 
using something like the docker stop command, or whatever else you prefer.

# Assumptions Made
In creating this program, some assumptions did have to be made.

## Strictness
It was assumed that there were very strict standards on the receipt format.
Beyond just having the data, it was throughly checked that the data followed
the proper formats when indicated, such as no negative prices, correct dates and times,
etc. Thus, there's a hefty amount of built in checking to make sure the receipt is
a valid one. 

## Time Ranges
The instructions clearly indicated that the valid time range was after 2 pm and before 4 pm 
to get the extra points for that part of the point calculations. 

This was interpreted to mean that times of 14:01 all the way to 15:59 inclusive were valid.
If this is too lenient and only times in the 3 pm zone work, then this will cause issues.
If not, great!

## Lack of Error Diversity
The API seemed to make it pretty clear that only two common errors needed to exist:
a BadRequest, which contains the text "The receipt is invalid." and 
NotFound which holds the text "No receipt found for that ID."

Because of this, only those two error texts and their appropriate HTTP statuses are shown
in most cases. What results is that any time a check fails for the receipt, BadRequest is thrown,
and when the database is queried with an ID that doesn't exist, NotFound is thrown.

There are a couple exceptions but those errors are very much not likely to be thrown.

## No Receipt Storage
The receipts are not stored in the database, since only points are ever queried. 
This can be easily changed but for now it's a practical act of saving a bit of 
data on the backend.

# Conclusion
Provided this application can run, it should be pretty much up to snuff 
with the requirements and demonstrate some decent Go code.

Enjoy looking at this and testing it!



