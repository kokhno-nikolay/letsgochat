var submit_btn = document.getElementById("form-submit")
var chat_inp = document.getElementById("chat_inp")
var socket = new WebSocket("ws://letsgochat.herokuapp.com/chat?token=23fddf65-f18f-4674-92a8-d8869bd13f7e");

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