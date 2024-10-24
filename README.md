**Introduction**
This is my Todo app for the GoLang Academy. 

**App Structure**
The app is comprised of:
- [cmd/server.go] A web server responsible for routing api URIs to an appropriate handler and hosting the web frontend.
- [cmd/web] The frontend web app. Simple web page that allows a user to create, mark as complete, and delete Todo items from a Todo list.
- [api/] The api connecting the web server to the data store.
- [data/datastore.go] Where the in-memory todo item list is contained.
- [services/dataService.go] A service used to manipulate the data within the data store. Called by the api.
- [utils] Just some reusable code for strings and slices.
