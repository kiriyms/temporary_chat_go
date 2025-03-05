## Chat web application
Made in Go + HTMX.

### Table of contents
- [Features](#features)
- [Demo](#demo)
- [Documentation](#documentation)
- [Technologies](#technologies)
- [Techniques & paradigms](#techniques--paradigms)
- [Reflection on challenges](#reflection-on-challenges)
- [Known issues](#known-issues)

##### **Features**
- Quick anonymous chatting
- Basic personalization functionality: upload avatar and change nickname
- Creating chat rooms and joining rooms via short numeric codes
- User session tracking via cookies

##### **Demo**
It is possible to demo the app at http://164.90.179.107:8080

The demo is hosted via a droplet on DigitalOcean.
Alternatively, it's possible to clone this repository, build and run the app locally.

- To build the app, type `go build -o main.exe ./cmd` (Windows) OR `go build -o main ./cmd` (Linux). 

"main.exe" or "main" will be the binary files with the server. Running those files will start the server at `http://localhost:1323`.
Additionally, a .env file needs to be added with a "JWT_SECRET" field, used to sign JWTokens for session and authorization tracking.

#### Documentation
Navigating the app happens in several steps:

1. Upon first entering, you will be prompted to enter a nickname and upload an avatar. This is not mandatory, leaving the fields blank will result in an "Anonymous" nickname and a default avatar picture.
2. After the initial page, you will be redirected to the "dashboard", which is the main space for using the app. Here you can do several things:
    1. Edit nickname/avatar
    2. Create a new chat room
    3. Join an existing chat room via numeric code (you need to ask another user for the code to their room)
    4. Navigate up to 5 different rooms that you are a part of
3. Within a room, you are able to see the participants, write and read messages. It is also possible to exit the room.

> **Note!**
>
> The rooms are limited by a 3-minute timer. This is done due to a technical decision explained in [Techniques & paradigms](#techniques--paradigms) section

#### Technologies
The application was done using the following technologies:

| Aspect | Technology |
|--------|------------|
|Backend API|Go ([Echo](https://echo.labstack.com) library)|
|Frontend interface|HTML+CSS ([HTMX](https://htmx.org) library for handling server reponses)|
|Session authorization and tracking|JWT and Cookies|

#### Techniques & paradigms
- **HATEOAS**: this app does not have a separate web client. Most of the interactivity and all of the responses are sent as HTML packages directly to the browser.
- **Dependency injection**: the handler for API endpoints is done using Go interfaces. This allowed for a creation of two easily-interchangeable handler instances which facilitated early manual testing of endpoints, before HTML rendering was set up.
- **Concurrency**: multiple rooms, their timers and websocket tickers are all run concurrently using Go's goroutine/channel systems.
- **In-memory data storage**: to keep the tech stack of the project simpler, I decided to avoid using any database solution to store data. As such, there is no persistent data in the application, and all user and room data is timed. Generally, the time limit is 3 minutes.

#### Reflection on challenges
- The HATEOAS does not easily support client-side interactivity, however, I wanted to try out setting up some simple interactive animations. To achieve that, I used HTMX lifecycle events to act on the server responses and hook up Javascript functions in the browser that provided some limited intercativity.
- By utilizing Go+HTMX stack, and Cookies to transfer JWT, I learned more about the structure of HTTP requests and responses in relation to headers, content-types, request/response bodies and other elements.
- Uploaded avatar pictures need to persist for the duration of the user session. However, as the user itself does not persist, there is no reason to keep the uploaded avatar on server disk storage. As a result, I set up a timer system to delete the image after the session is over.

#### Known issues
>At the moment, messages that are sent in chat by users are not filtered, which allows for XSS (Cross Site Scripting) attacks. 

If you'd like to try out the app, please make sure you only join rooms of your trusted friends, or simulate different users by opening several incognito browser tabs (such that the browser does not share the local Cookie storage)
