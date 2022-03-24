var submit_btn = document.getElementById("form-submit")
var chat_inp = document.getElementById("chat_inp")
var socket = new WebSocket("ws://letsgochat.herokuapp.com/chat?token=4e4e4a66-ab92-43a7-9a6a-44d5c67a1b3a");

socket.onopen = () => {
    console.log("Successfully connected from ws server");
};

socket.onclose = event => {
    console.log("Socket closed connection: ", event);
    socket.send("Client closed connection!")
};

socket.onerror = error => {
    console.log("Socket error: ", error);
    document.write('<html><body>invalid token</body></html>')
};

socket.onmessage = function(event) {
    var out = document.getElementById('res');
    out.innerHTML = event.data;
}

submit_btn.addEventListener('click', function () {
    sendMessage(chat_inp.value)
});

function sendMessage(data) {
    socket.send(data)
}