# Fetch Backend API - Receipt Processor
## Coded by Jesse Jones

# Building Docker Image
## To create the Docker image for this program:

 0. Be on a unix-like environment. If you're on Windows, use the Dockerfile 
    and Go files to build it using some different method.
 1. Change directory into `webserver` directory held in this repo.
 2. Run the bash script `buildImage.sh`, which will run the commands necessary 
    to generate a Docker image based on the Dockerfile and the three Go files provided in the directory.
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
    You can just do a find and replace wherever 8000 shows up.

# Testing the Server
Beyond the automatic testing that will be done, a few example jsons are provided with
a couple bash scripts that can insert them into the database via a curl request or
query the server with a given ID.

These examples and test scripts are found in the `tests` directory in the repo. 

The script `insertReceiptToDb` with an argument of a JSON file, is what's used
to more easily insert a JSON receipt into the database to then get an id
JSON back if it's a valid receipt.

The script `queryPointsAtID.sh` takes an ID that can be found from the output
of the server, and produces the point count associated with it if it's a valid ID.

# Stopping the Server
To stop the server, just stop the container from running via killing it directly, 
using something like the `docker stop` command, or whatever else you prefer.

# Assumptions Made
In creating this program, some assumptions did have to be made.

## Strictness
It was assumed that there were very strict standards on the receipt format.
Beyond just having the data, it was throughly checked that the data followed
the proper formats when indicated, such as no negative prices, correct dates and times,
etc. Thus, there's a hefty amount of built-in checking to make sure the receipt is
a valid one. 

## Time Ranges
The instructions clearly indicated that the valid time range was after 2 pm and before 4 pm 
to get the extra points for that part of the point calculations. 

This was interpreted to mean that times of 14:01 all the way to 15:59 inclusive were valid.
If this is too lenient and only times in the 15 hour zone work, then this will cause issues.
If not, great! It's a fairly reasonable assumption to make since 14:01 to 15:59 fulfils the
after 2 pm and before 4 pm criterion.

## Lack of Error Diversity
The API seemed to make it pretty clear that only two common errors needed to exist:
a BadRequest, which contains the text "The receipt is invalid." and 
NotFound which holds the text "No receipt found for that ID."

Because of this, only those two error texts and their appropriate HTTP statuses are shown
in most cases. What results is that any time a check fails for the receipt, BadRequest is thrown,
and when the database is queried with an ID that doesn't exist, NotFound is thrown.

There are a couple exceptions but those errors are very much 
not likely to be thrown and thus can be ignored for the general case, since both of them 
are catching errors for encrypting response structs to JSON. Thus, if they do happen,
they are labled as internal server errors since them breaking isn't the user's fault. 

## No Receipt Storage
The receipts are not stored in the database, since only points are ever queried. 
This can be easily changed but for now it's a practical act of saving a bit of 
data on the backend. Plus, the databse itself is only a memory-based data-structure,
so it's not made for longterm use anyway.

## Needing to Deal with Multithreading.
The database struct not only contaiins a hashmap but also a mutex, meaning that 
any one thread needs to lock down the mutex first before altering/querying the database,
unlocking it after use. This ensures that there are no simultaenious databse changes, allowing
a sort of queue to form if a ton of queries are made at once. This was done since each endpoint 
function does run in its own seperate quasi-thread type thing called a goroutine. Thus, a mutex
ensures no overlaping edits or queries. Is this necessary for the scope of this kind of assignment?
Probably not. Just good to have anyway.

## Locality of Program
It was also assumed that the server will only run over 
a local network on one machine with the queries happening on the same machine.
Some tweaking can be made to query from a different machine on the same network 
but the local IP address of the server device will need to be known.

To make this application able to work across the Internet and not just over LAN,
a port will need to be forwarded and the public ip will be needed along with the 
forwarded port. 

Security-wise this is a nightmare so it's probably just a local webserver, hence
why the assumption of it being local was made.

# Conclusion
Provided this application can run on your machine 
using your installation of Docker, 
it should be pretty much up to snuff 
with the requirements and demonstrate some decent Go code.

Enjoy looking at this and testing it!



