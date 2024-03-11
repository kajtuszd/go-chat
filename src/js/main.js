let username = ''

const joinForm = document.getElementById('join-form')
joinForm.addEventListener('submit', event => {
    event.preventDefault()
    username = document.getElementById("username-input").value.trim()
    if (username !== '') {
        joinForm.style.display = 'none'
        connectWebSocket()
    }
})

let socket

function connectWebSocket() {
    socket = new WebSocket(`ws://localhost:8080/ws?username=${username}`)

    socket.addEventListener('open', () => {
        socket.send(username + ' joined the chat')
        console.log('WebSocket connection established')
    })

    socket.addEventListener('message', event => {
        displayMessage(event.data)
    })

    socket.addEventListener('close', () => {
        displayMessage(username + ' left the chat')
        console.log('WebSocket connection closed')
        socket.close()
    })
}

document.getElementById("close-button").addEventListener("click", () => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(username + " left the chat")
        socket.close()
    } else {
        console.log("WebSocket connection is not open.")
    }
})

window.addEventListener("beforeunload", () => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(username + " exited the chat")
        socket.close()
    }
})

function displayMessage(message) {
    let chatWindow = document.getElementById("chat-window")
    chatWindow.innerHTML += message + "</br>"
    chatWindow.scrollTop = chatWindow.scrollHeight
}

function sendMessage() {
    let message = document.getElementById("message-input").value.trim()
    const content = `${username} : ${message}`
    if (content !== '') {
        socket.send(content)
        message.value = ''
    }
}

function toggleForms() {
    const joinForm = document.getElementById("join-form")
    const chatForm = document.getElementById("chat-form")

    if (joinForm.style.display === "none") {
        joinForm.style.display = "block"
        chatForm.style.display = "none"
    } else {
        joinForm.style.display = "none"
        chatForm.style.display = "block"
    }
}

function checkJoinForm() {
    const roomCheckbox = document.getElementById("room-checkbox");
    const joinButton = document.getElementById("join-button")
    const roomInput = document.getElementById("room-input");

    const usernameInput = document.getElementById("username-input")

    if (roomCheckbox.checked) {
        roomInput.style.display = "block";
    } else {
        roomInput.style.display = "none"
    }

    if (usernameInput.value.trim() === "") {
        joinButton.disabled = true
    } else if (roomCheckbox.checked && roomInput.value.trim() === "") {
        joinButton.disabled = true
    } else {
        joinButton.disabled = false
    }
}

function checkMessage() {
    const messageInput = document.getElementById("message-input")
    const messageButton = document.getElementById("message-button")

    messageButton.disabled = messageInput.value.trim() === ""
}
