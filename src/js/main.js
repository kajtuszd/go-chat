const socket = new WebSocket("ws://localhost:8080/ws")

socket.addEventListener('open', (event) => {
    socket.send('New Anon joined!');
});

socket.addEventListener('message', (event) => {
    console.log('Message from server: ', event.data);
    displayMessage(event.data)
});

function sendMessage() {
    let message = document.getElementById("message-text")
    let messageText = message.value
    message.value = ""
    socket.send(messageText)
}

function displayMessage(message) {
    let chatWindow = document.getElementById("chat-window")
    chatWindow.innerHTML += 'Anonymous: ' + message + '</br>'
    chatWindow.scrollTop = chatWindow.scrollHeight
    return false
}
