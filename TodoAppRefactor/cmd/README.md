## Frontend Web App

Contained within the 'web' folder, the frontend of the app is a basic web page that allows a user to:
- Add new todo items
- Mark todo items as complete
- Delete todo items

The frontend calls the API from 'onclick' commands on the corresponding buttons. When the page is originally loaded, it calls
the 'getAll' API call to retrieve and then display any existing todo items.

## Server

The server is responsible for a couple of things:
- Sets up a file server to serve static files such as the stylesheet and an image.
- Serve the frontend web page.
- Sets up the API 'RequestHandler' as well as a stop channel that is used to shut down the Request handler.
- Sets up the API routes to the corresponding handlers
- Attempts to "gracefully" shut down the request handler but after running the frontend, it appears this code isn't executed
  when pressing 'ctrl+c'. The print statements are not printed.
