var submit_btn = document.getElementById("form-submit")
var chat_inp = document.getElementById("chat_inp")
var socket = new WebSocket("ws://letsgochat.herokuapp.com/chat?token=6bda2cbf-b778-4ae2-b7d4-354f0616c94b");

socket.onopen = () => {
    console.log("Successfully connected from ws server");
};

socket.onclose = event => {
    console.log("Socket closed connection: ", event);
    socket.send("Client closed connection!")
};

socket.onerror = error => {
    console.log("Socket error: ", error);
    sdocument.write('<html><body>invalid token</body></html>')
};

socket.onmessage = function(event) {
    let data = JSON.parse(event.data);
    let out = document.getElementById('messages');
    let chatContent = `<p><strong>${data.username}</strong>: ${data.text}</p>`;
    out.innerHTML += chatContent
}

submit_btn.addEventListener('click', function () {
    sendMessage(chat_inp.value)
});

function sendMessage(data) {
    socket.send(data)
}