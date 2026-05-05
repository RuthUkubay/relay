export interface User {
  id: string;
  name: string;
}

export interface Room {
  id: string;
  name: string;
}

export interface Message {
  id: string;
  userId: string;
  username: string;
  content: string;
  room: string;
  timestamp: string;
}

export async function createUser(name: string): Promise<User> {
  const res = await fetch('/api/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<User>;
}

export async function fetchRooms(): Promise<Room[]> {
  const res = await fetch('/api/rooms');
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<Room[]>;
}

export async function fetchMessages(room: string): Promise<Message[]> {
  const res = await fetch(`/api/messages?room=${encodeURIComponent(room)}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<Message[]>;
}

export async function sendMessage(userId: string, content: string, room: string): Promise<Message> {
  const res = await fetch('/api/messages', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ userId, content, room }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json() as Promise<Message>;
}
