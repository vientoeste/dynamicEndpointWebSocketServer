import { WebSocketServer } from "ws";

const wss = new WebSocketServer({
  port: 3000,
});
const endpoints = {
  new: [],
  applied: [],
};

wss.on('connection', (socket) => {
  socket.on('message', (message) => {
    const data = JSON.parse(message);
    if (typeof data === 'string') {
      throw new Error('received string data');
    }

    if (data.action === 'addRoute') {
      endpoints.new.push(data.route);
      updateServerBehavior();
      wss.clients.forEach((client) => {
        // [TODO] update room lists here
        client.send(JSON.stringify({ action: 'updateRoutes', routes: endpoints }));
      });
    }
    wss.clients.forEach((cli) => {
      cli.send(JSON.stringify({
        action: 'updateRoutes',
        route: data.route,
      }));
    })
  })
});

const updateServerBehavior = () => {
  console.log(1232131)
  endpoints.new.forEach((endpoint) => {
    const hasRouteHandler = wss.eventNames().includes('connection:' + endpoint);
    if (!hasRouteHandler) {
      console.log(12321432)
      wss.on('connection:' + endpoint, (socket) => {
        // [TODO] Handle messages for this route
        socket.send(123123)
        console.log(37)
      });
    }
  });
};
