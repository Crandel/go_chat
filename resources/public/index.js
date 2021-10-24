let sock = new WebSocket('ws://' + window.location.host + '/ws');

let colors = ['green', 'red', 'yellow', 'blue'],
    fillCollor = Object.create(null);

sock.onopen = () => {
  console.log("Open web socket");
};
// income message handler
sock.onmessage = (event) => {
  let messageElem = document.getElementById('canvas')
  let obj = JSON.parse(event.data);
  let span = document.createElement("span")
  let br = document.createElement("br")
  span.setAttribute("class", "message")
  span.append(obj.message)
  span.append(br)
  messageElem.append(span)
  console.log("New event");
  console.log(obj);
};

sock.onclose = (event) => {
  console.log("Connection closed");
  console.log(event);
};

sock.onerror = (error) => {
  console.log("Error:", error);
};

sendMessage = (obj) => {
  sock.send(JSON.stringify(obj))
}

submitButton = () => {
  let message = document.getElementById("message").value
  if (message.length > 0){
    let obj = {'message': message};
    sendMessage(obj)
  }
}
