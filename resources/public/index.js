try{
  let sock = new WebSocket('ws://' + window.location.host + '/ws');
}
catch(err){
  let sock = new WebSocket('wss://' + window.location.host + '/ws');
}

let colors = ['green', 'red', 'yellow', 'blue']
