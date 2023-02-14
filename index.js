import { WebSocketServer } from "ws";

const wss = new WebSocketServer({
  port: 3000,
});
const endpoints = [];

wss.on('connection', (socket) => {
  socket.on('message', (message) => {
    const data = JSON.parse(message);
    if (typeof data === 'string') {
      throw new Error('received string data');
    }

    if (data.action === 'addRoute') {
      endpoints.push(data.route);
    }
    wss.clients.forEach((cli) => {
      cli.send(JSON.stringify({
        action: 'updateRoutes',
        route: data.route,
      }));
    })
  })
});

endpoints.forEach((route) => {
  wss.options('connection', (socket) => {
    if (socket.route === route) {

    }
  })
});