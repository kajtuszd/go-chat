const socket = new WebSocket("ws://localhost:8080/ws")

socket.addEventListener('open', (event) => {
    socket.send('New Anon joined!');
});

socket.addEventListener('message', (event) => {
    // console.log('Message from server: ', event.data)
    displayMessage(event.data)
});

socket.addEventListener('close', (event) => {
    console.log('Connection closed by exit button: ', event.data)
    // displayMessage(event.data)
});

socket.addEventListener('error', (event) => {
    console.log('An error occurred: ', event.data)
    displayMessage(event.data)
});

document.getElementById("closeButton").addEventListener('click', () => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send('Anon left the chat')
        socket.close();
    } else {
        console.log('WebSocket connection is not open.');
    }
});

window.addEventListener('beforeunload', () => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send("Anon exited the chat")
        socket.close();
    }
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
