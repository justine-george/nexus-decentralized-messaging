let socket;
let rtcPeerConnection;
let dataChannel;

function connect() {
  socket = new WebSocket("wss://" + window.location.host + "/ws");

  socket.onopen = function (e) {
    console.log("WebSocket connection established");
  };

  socket.onmessage = function (event) {
    const message = JSON.parse(event.data);
    handleMessage(message);
  };

  socket.onclose = function (event) {
    console.log("WebSocket connection closed");
  };
}

function handleMessage(message) {
  switch (message.type) {
    case "peer_list":
      updatePeerList(JSON.parse(message.content));
      break;
    case "offer":
      handleOffer(message);
      break;
    case "answer":
      handleAnswer(message);
      break;
    case "ice_candidate":
      handleIceCandidate(message);
      break;
    case "chat":
      displayMessage(message);
      break;
  }
}

function updatePeerList(peers) {
  const peerList = document.getElementById("peers");
  peerList.innerHTML = "";
  for (const [id, name] of Object.entries(peers)) {
    const peerElement = document.createElement("div");
    peerElement.textContent = name;
    peerElement.onclick = () => connectToPeer(id);
    peerList.appendChild(peerElement);
  }
}

function connectToPeer(peerId) {
  // Implement WebRTC connection logic here
}

function sendMessage() {
  const input = document.getElementById("messageInput");
  const message = input.value;
  input.value = "";

  if (dataChannel && dataChannel.readyState === "open") {
    dataChannel.send(JSON.stringify({ type: "chat", content: message }));
  } else {
    socket.send(JSON.stringify({ type: "chat", content: message }));
  }

  displayMessage({ from: "You", content: message });
}

function displayMessage(message) {
  const messagesDiv = document.getElementById("messages");
  const messageElement = document.createElement("div");
  messageElement.textContent = `${message.from}: ${message.content}`;
  messagesDiv.appendChild(messageElement);
}

document.getElementById("sendButton").onclick = sendMessage;

connect();
