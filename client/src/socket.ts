import type { Message, User } from './api';

export interface WSMessage {
  type: 'message' | 'user_join';
  message?: Message;
  user?: User;
}

type MessageHandler = (msg: WSMessage) => void;

// Opens a new WebSocket scoped to a specific room.
// The server uses the room query param to decide which broadcasts to forward.
export function connectWebSocket(userId: string, room: string, onMessage: MessageHandler): WebSocket {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const url = `${protocol}//${window.location.host}/ws?userId=${userId}&room=${encodeURIComponent(room)}`;
  const ws = new WebSocket(url);

  ws.onopen = () => console.log(`WebSocket connected (room: ${room})`);
  ws.onclose = () => console.log(`WebSocket disconnected (room: ${room})`);
  ws.onerror = (e) => console.error('WebSocket error', e);

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data as string) as WSMessage;
      onMessage(msg);
    } catch {
      console.error('Failed to parse WebSocket message', event.data);
    }
  };

  return ws;
}
