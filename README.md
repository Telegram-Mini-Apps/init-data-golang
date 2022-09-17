# Telegram Web Apps Init Data

Telegram Web Apps init data is rather important part of Telegram`s platform. 
This repository contains its explanation and examples verification via 
different languages. You can find official verification example via 
pseudocode in Telegram Web Apps [documentation](https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app).

## Table of contents
- [Core information](#core-info)
  - [How parameters are passed into app](#how-params-are-passed)
  - [Passing parameters to server](#passing-to-server)
- [Verification examples](#verification-examples)
- [Contribution](#contribution)

<a name="core-info"/>

## Core information

Launching an application on Telegram Web Apps platform implies passing 
special parameters which are connected with current user. In this 
documentation these parameters (or _initial data_) are called 
**launch parameters**.

According to TWA applications are usual web applications, and they should be
correctly displayed on any device, they are always wrapped into WebView. So,
controlling device can communicate with our web application through WebView
functionality.

<a name="how-params-are-passed"/>

### How parameters are passed into app

The easiest way to pass launch parameters to application and allow application
to get them while launching javascript code is to mention them in application
URL. That's why Telegram Web Apps uses this way.

According to [official documentation](https://core.telegram.org/bots/webapps#initializing-web-apps),
it is required to add special script into `<head>` tag. This action leads 
to process which creates a new field in `window` object. As a result, you will
be able to communicate with native device through `window.Telegram` object.

In this documentation it is required to operate only with 
`window.Telegram.WebApp` object:

![img.png](assets/webapp-window-obj.png)

<a name="passing-to-server"/>

### Passing parameters to server

One of the main features of launch parameters is they could be used as 
authorization factor. It means, you could use them to identify requesting
client.

As long as launch parameters are always signed by Telegram bot secret 
token (sign is placed in `hash` parameter), you could always verify them and
trust passed parameters.

The best way to pass launch parameters to your server is to specify them in
some header. Here comes an example in JavaScript with `axios` library usage:

```javascript
import axios from 'axios';

// Create axios instance.
const http = axios.create({
  headers: {
    // Append authorization header.
    Authorization: `Bearer ${window.Telegram.WebApps.initData}`,
  }
});

// Now, in case we use this instance to perform requests, authorization header
// will be automatically appended. The next thing we have to do is just verify 
// it on server side.
```

<a name="verification-examples"/>

## Verification examples

- [Golang](/examples/golang)
- [Node JS](/examples/node)

## Contribution

Any contribution is appreciated. Please, use already existed examples to 
create new one.
