HOW TO START THE SERVICE:
1. Open up your terminal and navigate to the path (or use the terminal in the IDE)
2. Run the command "go run *.go"
3. You should see "Starting Application" in the terminal to know the program has compiled and started
   <br />

GENERAL POSTMAN INFO:
1. I have defined the port to be "12345"
   <br />
APIS:
I have created 5 APIs to be used from.
1. Create User (Insert a new user into the database if user does not exist)
2. Update User (Update user's properties (except user_id))
3. Get User (Get a single user and all of the user's properties)
4. Get All Users (Get a list of all the users in the database)
5. Delete User (Remove a single user from the database)

CURL REQUESTS:
1. CREATE USER:<br/>
   curl --location --request POST 'http://localhost:12345/user' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "user_id":"1",
   "name": "user1",
   "dob": {
   "day" : 1,
   "month" : 10,
   "year" : 1996
   },
   "address": {
   "street": "Street",
   "block" : "123",
   "unit": "123"
   },
   "description": "description1"
}'
<br /><br />

2. UPDATE USER:<br />
   curl --location --request PATCH 'http://localhost:12345/user/1' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "name": "updateduser1",
   "dob": {
   "day" : 10,
   "month" : 12,
   "year" : 2000
   },
   "address": {
   "street": "UpdatedStreet",
   "block" : "Updated123",
   "unit": "Updated123"
   }
   }'
<br /><br />

3. GET USER:<br />
   curl --location --request GET 'http://localhost:12345/user/1'
<br /><br />
   
4. GET ALL USERS:<br />
   curl --location --request GET 'http://localhost:12345/user'
<br /><br />
   
5. DELETE USER:<br />
   curl --location --request DELETE 'http://localhost:12345/user/1'


<br />

EXTRA INFO:
1. Code written using Golang & MongoDB
2. I did not write unit tests, attempted but failed due to not knowing how to mock the mongodb package (tried for quite a long time)
3. Did not attempt the advanced requirements, but happy, open and willing to discuss how to implement them
