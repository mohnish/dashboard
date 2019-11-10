let ws = new WebSocket("ws://" + window.location.host + "/ws");

ws.onopen = () => {
  console.log('connected');
};

ws.onmessage = (evt) => {
  const widget = JSON.parse(evt.data);
  console.log("Received message", widget);
  if (widget.id) {
    updateMap(widget)
    renderWidgets();
  }
}

ws.onclose = () => {
  console.log('your browser window will automatically be refreshed');
  setTimeout(() => location.reload(true), 3e3);
}
