import { createUser, fetchRooms, fetchMessages, sendMessage, type User } from './api';
import { connectWebSocket, type WSMessage } from './socket';
import { renderMessage, clearMessages, renderRooms, setActiveRoom, showChat, setUserLabel } from './ui';

let currentUser: User | null = null;
let currentRoom: string | null = null;
let ws: WebSocket | null = null;

// Switches to a room: reconnects the WebSocket, fetches history, re-renders.
async function switchRoom(roomName: string): Promise<void> {
  if (currentRoom === roomName) return;
  currentRoom = roomName;

  if (ws) ws.close(1000, "switching rooms");
  ws = connectWebSocket(currentUser!.id, roomName, (msg: WSMessage) => {
    if (msg.type === 'message' && msg.message) {
      renderMessage(msg.message);
    }
  });

  clearMessages();
  setActiveRoom(roomName);

  const history = await fetchMessages(roomName);
  history.forEach(renderMessage);
}

async function handleLogin(name: string): Promise<void> {
  currentUser = await createUser(name);
  setUserLabel(currentUser.name);

  const rooms = await fetchRooms();
  const firstRoom = rooms[0]?.name ?? 'general';

  renderRooms(rooms, firstRoom, (name) => {
    switchRoom(name).catch(console.error);
  });

  showChat();
  await switchRoom(firstRoom);
}

async function handleSend(content: string): Promise<void> {
  if (!currentUser || !currentRoom || !content.trim()) return;
  // Server saves, then broadcasts back via WebSocket — we don't render here.
  await sendMessage(currentUser.id, content.trim(), currentRoom);
}

document.addEventListener('DOMContentLoaded', () => {
  const loginForm = document.getElementById('login-form') as HTMLFormElement;
  const loginInput = document.getElementById('login-input') as HTMLInputElement;
  const messageForm = document.getElementById('message-form') as HTMLFormElement;
  const messageInput = document.getElementById('message-input') as HTMLInputElement;

  loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const name = loginInput.value.trim();
    if (name) handleLogin(name).catch(console.error);
  });

  messageForm.addEventListener('submit', (e) => {
    e.preventDefault();
    handleSend(messageInput.value).catch(console.error);
    messageInput.value = '';
  });
});
