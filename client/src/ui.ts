import type { Message, Room } from './api';

export function renderMessage(msg: Message): void {
  const list = document.getElementById('messages')!;
  const item = document.createElement('li');
  const time = new Date(msg.timestamp).toLocaleTimeString();
  item.textContent = `[${time}] ${msg.username}: ${msg.content}`;
  list.appendChild(item);
  list.scrollTop = list.scrollHeight;
}

export function clearMessages(): void {
  document.getElementById('messages')!.innerHTML = '';
}

export function renderRooms(rooms: Room[], activeRoom: string, onSelect: (name: string) => void): void {
  const list = document.getElementById('rooms-list')!;
  list.innerHTML = '';
  for (const room of rooms) {
    const li = document.createElement('li');
    li.textContent = `# ${room.name}`;
    li.dataset['room'] = room.name;
    if (room.name === activeRoom) li.classList.add('active');
    li.addEventListener('click', () => onSelect(room.name));
    list.appendChild(li);
  }
}

export function setActiveRoom(name: string): void {
  document.querySelectorAll('#rooms-list li').forEach((el) => {
    el.classList.toggle('active', (el as HTMLElement).dataset['room'] === name);
  });
  (document.getElementById('room-header') as HTMLElement).textContent = `# ${name}`;
}

export function showChat(): void {
  (document.getElementById('login-section') as HTMLElement).style.display = 'none';
  (document.getElementById('chat-section') as HTMLElement).style.display = 'flex';
}

export function setUserLabel(name: string): void {
  (document.getElementById('user-label') as HTMLElement).textContent = name;
}
