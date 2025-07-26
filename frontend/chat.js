let socket;
console.log("js is loaded")
function initChat() {
  if (socket && socket.readyState === WebSocket.OPEN) return;

  const messagesDiv = document.getElementById('messages');
  const inputTo     = document.getElementById('chat-to');
  const inputText   = document.getElementById('chat-text');
  const btnSend     = document.getElementById('chat-send');
  if (!messagesDiv || !inputTo || !inputText || !btnSend) return;

  const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
  const wsUrl = proto + location.host + '/api/chat/ws';
  console.log('Connecting WS to', wsUrl);

  socket = new WebSocket(wsUrl);

  socket.onopen = () => console.log('WS connected');

  socket.onmessage = ev => {
     console.log('Raw WS message:', ev.data); // 👈
  try {
    const msg = JSON.parse(ev.data);
    console.log('Received:', msg);
    
    const sender = msg.sender || '???';
    const receiver = msg.receiver || '???';
    const text = msg.text || '(no text)';
    const timestamp = msg.timestamp ? new Date(msg.timestamp).toLocaleTimeString() : '??:??:??';

    const line = document.createElement('div');
    line.textContent = `[${timestamp}] ${sender} → ${receiver}: ${text}`;
    messagesDiv.appendChild(line);
  } catch (e) {
    console.error('Invalid message', e);
  }
};

  socket.onclose = () => console.log('WS disconnected');
  socket.onerror = e => console.error('WS error', e);

  btnSend.onclick = () => {
    const receiver = inputTo.value.trim();
    const text = inputText.value.trim();

    if (!receiver || !text) return;

    console.log('Sending message:', { receiver, text });

    if (socket.readyState !== WebSocket.OPEN) {
      alert('Socket not open');
      return;
    }

    const msg = {
      receiver: receiver,
      text: text
    };

    socket.send(JSON.stringify(msg));
    inputText.value = '';
  };
}
initChat();
