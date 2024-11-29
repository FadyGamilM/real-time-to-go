const WebSocket = require('ws');

// Create a WebSocket connection
const socket = new WebSocket("ws://localhost:3000/ws");

// Listen for messages from the server
socket.on("message", (data) => {
  console.log("Received from the server: ", data.toString());
});

// Send a message to the server once the connection is open
socket.on("open", () => {
  console.log("Connection established with the server");
  socket.send("Hello from client");
});

// Handle any errors that occur
socket.on("error", (error) => {
  console.error("WebSocket error: ", error);
});

// Log when the connection is closed
socket.on("close", () => {
  console.log("Connection to the server closed");
});
