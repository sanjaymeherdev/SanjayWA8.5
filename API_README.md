# GoWA API Reference
## WhatsApp Gateway — v8.5.0 (Multi-Tenant)

**Base URL (production):** `https://graphicy-wa--testhimrealty.replit.app`  
**Base URL (dev):** `https://<your-replit-dev-domain>`

All requests that target a specific WhatsApp account must include the device ID header:

```
X-Device-Id: <device_id>
```

If only **one** device is registered, the header is optional — the server uses it automatically.

---

## Table of Contents

1. [Device Management](#1-device-management)
2. [Account / App Control](#2-account--app-control)
3. [Send Messages](#3-send-messages)
4. [Chats](#4-chats)
5. [Users](#5-users)
6. [Groups](#6-groups)
7. [Messages (Actions)](#7-messages-actions)
8. [AI Chat (Groq)](#8-ai-chat-groq)
9. [WebSocket Events](#9-websocket-events)
10. [Apps Script Integration Guide](#10-apps-script-integration-guide)
11. [Common Response Format](#11-common-response-format)

---

## 1. Device Management

These endpoints manage which WhatsApp accounts (devices) are connected to the server. **No `X-Device-Id` header needed** for these.

### List all connected devices
```
GET /devices
```
**Response:**
```json
{
  "status": 200,
  "code": "SUCCESS",
  "results": [
    {
      "id": "628123456789@s.whatsapp.net",
      "name": "My Phone",
      "platform": "android",
      "is_connected": true,
      "is_logged_in": true,
      "last_seen": "2026-05-07T12:00:00Z"
    }
  ]
}
```

### Add a new device slot
```
POST /devices
```
**Body (JSON):**
```json
{ "name": "Client A Phone" }
```
**Response:** Returns `device_id` to use in subsequent calls.

### Get a specific device
```
GET /devices/:device_id
```

### Show QR code for login
```
GET /devices/:device_id/login
```
Opens a page with a QR code to scan in WhatsApp → Linked Devices.

### Login with pairing code (no QR scan)
```
POST /devices/:device_id/login/code
```
**Body (JSON):**
```json
{ "phone": "628123456789" }
```
Returns an 8-character pairing code to enter in WhatsApp → Linked Devices → Link with phone number.

### Logout a device
```
POST /devices/:device_id/logout
```

### Reconnect a device
```
POST /devices/:device_id/reconnect
```

### Get device connection status
```
GET /devices/:device_id/status
```
**Response:**
```json
{
  "status": 200,
  "results": {
    "is_connected": true,
    "is_logged_in": true,
    "device_id": "628123456789@s.whatsapp.net"
  }
}
```

### Remove a device
```
DELETE /devices/:device_id
```

---

## 2. Account / App Control

> All endpoints below require `X-Device-Id` header (or auto-resolved if only one device).

### Legacy login (single device QR)
```
GET /app/login
```

### Check connection status (legacy)
```
GET /app/status
```

### List linked sub-devices
```
GET /app/devices
```

### Logout (legacy)
```
GET /app/logout
```

### Reconnect (legacy)
```
GET /app/reconnect
```

### Health check
```
GET /health
```
Returns `200 OK` with body `OK` when server is healthy. No auth required. Use this to ping the server from Apps Script.

---

## 3. Send Messages

> All send endpoints require `X-Device-Id` header. Use `multipart/form-data` for file uploads, `application/json` for text/JSON payloads.

### Send text message
```
POST /send/message
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "message": "Hello from GoWA!",
  "reply_message_id": ""
}
```
- `phone`: International format without `+`, e.g. `628123456789`
- `reply_message_id`: Optional — ID of message to quote/reply

---

### Send image
```
POST /send/image
```
**Body (multipart/form-data):**

| Field | Type | Description |
|---|---|---|
| `phone` | string | Recipient phone number |
| `caption` | string | Optional image caption |
| `image` | file | Image file (JPG/PNG/GIF) |
| `url` | string | OR use a public image URL instead of file |
| `compress` | boolean | Compress before sending (default: false) |
| `reply_message_id` | string | Optional reply target |

**Body (JSON with URL):**
```json
{
  "phone": "628123456789",
  "caption": "Check this out!",
  "url": "https://example.com/image.jpg",
  "compress": false
}
```

---

### Send video
```
POST /send/video
```
**Body (multipart/form-data or JSON):**

| Field | Type | Description |
|---|---|---|
| `phone` | string | Recipient |
| `caption` | string | Video caption |
| `video` | file | Video file (MP4) |
| `url` | string | OR public video URL |
| `compress` | boolean | Compress before sending |
| `gif_playback` | boolean | Send as looping GIF (v8.4+) |

---

### Send audio
```
POST /send/audio
```
**Body (multipart/form-data or JSON):**

| Field | Type | Description |
|---|---|---|
| `phone` | string | Recipient |
| `audio` | file | Audio file (MP3/OGG/M4A) |
| `url` | string | OR public audio URL |

---

### Send document/file
```
POST /send/file
```
**Body (multipart/form-data or JSON):**

| Field | Type | Description |
|---|---|---|
| `phone` | string | Recipient |
| `caption` | string | File caption |
| `file` | file | Any file |
| `url` | string | OR public file URL |

---

### Send location
```
POST /send/location
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "latitude": -6.2146,
  "longitude": 106.8451,
  "address": "Jakarta, Indonesia"
}
```

---

### Send contact card
```
POST /send/contact
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "contact_name": "John Doe",
  "contact_phone": "628987654321"
}
```

---

### Send link with preview
```
POST /send/link
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "link": "https://example.com",
  "caption": "Check this link"
}
```

---

### Send sticker
```
POST /send/sticker
```
**Body (multipart/form-data or JSON):**

| Field | Type | Description |
|---|---|---|
| `phone` | string | Recipient |
| `sticker` | file | Image (auto-converted to WebP) |
| `url` | string | OR public image URL |

---

### Send poll
```
POST /send/poll
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "question": "Favorite color?",
  "options": ["Red", "Blue", "Green"],
  "max_answer": 1
}
```

---

### Set presence (online/offline)
```
POST /send/presence
```
**Body (JSON):**
```json
{ "presence": "available" }
```
Values: `available`, `unavailable`

---

### Send chat presence (typing indicator)
```
POST /send/chat-presence
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "presence": "composing",
  "media": "text"
}
```

---

## 4. Chats

### List all chats
```
GET /chats
```
**Query params:**
- `limit` — max results (default 20)
- `after_date` — filter after date (ISO 8601)

**Response:**
```json
{
  "results": [
    {
      "jid": "628123456789@s.whatsapp.net",
      "name": "John Doe",
      "last_message": "Hello",
      "last_message_time": "2026-05-07T10:00:00Z",
      "unread_count": 2
    }
  ]
}
```

### Get messages in a chat
```
GET /chat/:chat_jid/messages
```
**Query params:**
- `limit` — max messages
- `before_date` — ISO date for pagination

### Pin/unpin chat
```
POST /chat/:chat_jid/pin
```
**Body (JSON):**
```json
{ "pin": true }
```

### Set disappearing timer
```
POST /chat/:chat_jid/disappearing
```
**Body (JSON):**
```json
{ "timer": 86400 }
```
Timer in seconds (0 = off, 86400 = 24h, 604800 = 7 days).

### Archive/unarchive chat
```
POST /chat/:chat_jid/archive
```
**Body (JSON):**
```json
{ "archive": true }
```

---

## 5. Users

### Get user info
```
GET /user/info?phone=628123456789
```

### Check if number has WhatsApp
```
GET /user/check?phone=628123456789
```
**Response:**
```json
{
  "results": {
    "is_registered": true,
    "jid": "628123456789@s.whatsapp.net"
  }
}
```

### Get user avatar
```
GET /user/avatar?phone=628123456789
```

### Change own avatar
```
POST /user/avatar
```
**Body (multipart):** `avatar` file field.

### Change display name
```
POST /user/pushname
```
**Body (JSON):**
```json
{ "name": "My Bot Name" }
```

### Get privacy settings
```
GET /user/my/privacy
```

### List my groups
```
GET /user/my/groups
```

### List my newsletters
```
GET /user/my/newsletters
```

### List my contacts
```
GET /user/my/contacts
```

### Get business profile
```
GET /user/business-profile?phone=628123456789
```

---

## 6. Groups

### Create group
```
POST /group
```
**Body (JSON):**
```json
{
  "title": "My Group",
  "participants": ["628111111111", "628222222222"]
}
```

### Join group with invite link
```
POST /group/join-with-link
```
**Body (JSON):**
```json
{ "link": "https://chat.whatsapp.com/AbcXyz123" }
```

### Get group info
```
GET /group/info?group_id=12345678@g.us
```

### Leave group
```
POST /group/leave
```
**Body (JSON):**
```json
{ "group_id": "12345678@g.us" }
```

### List participants
```
GET /group/participants?group_id=12345678@g.us
```

### Add participants
```
POST /group/participants
```
**Body (JSON):**
```json
{
  "group_id": "12345678@g.us",
  "participants": ["628333333333"]
}
```

### Remove / Promote / Demote participants
```
POST /group/participants/remove
POST /group/participants/promote
POST /group/participants/demote
```
Same body structure as Add participants.

### Get group invite link
```
GET /group/invite-link?group_id=12345678@g.us
```

### Set group name / topic / photo / lock / announce
```
POST /group/name
POST /group/topic
POST /group/photo
POST /group/locked
POST /group/announce
```

---

## 7. Messages (Actions)

### React to message
```
POST /message/:message_id/reaction
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "emoji": "👍"
}
```

### Revoke (delete for everyone)
```
POST /message/:message_id/revoke
```
**Body (JSON):**
```json
{ "phone": "628123456789" }
```

### Delete (for me only)
```
POST /message/:message_id/delete
```

### Edit message
```
POST /message/:message_id/update
```
**Body (JSON):**
```json
{
  "phone": "628123456789",
  "message": "Updated text"
}
```

### Mark as read
```
POST /message/:message_id/read
```
**Body (JSON):**
```json
{ "phone": "628123456789" }
```

### Star / Unstar message
```
POST /message/:message_id/star
POST /message/:message_id/unstar
```

### Download media from message
```
GET /message/:message_id/download?phone=628123456789
```

---

## 8. AI Chat (Groq)

Powered by the Groq API (LLaMA 3.3 70B). **No `X-Device-Id` required.**

### Chat with AI
```
POST /api/ai/chat
```
**Body (JSON):**
```json
{
  "message": "Summarize this WhatsApp conversation",
  "system_prompt": "You are a helpful WhatsApp assistant. Be concise.",
  "conversation_history": [
    { "role": "user", "content": "Previous message" },
    { "role": "assistant", "content": "Previous reply" }
  ]
}
```
**Response:**
```json
{
  "status": 200,
  "code": "SUCCESS",
  "results": {
    "response": "Here is a summary...",
    "model": "llama-3.3-70b-versatile"
  }
}
```
**Environment variable:** `GROQ_API_KEY` must be set. Model can be changed via `GROQ_MODEL` env var.

---

## 9. WebSocket Events

Connect to receive real-time events for a specific device:
```
wss://<host>/ws?device_id=628123456789@s.whatsapp.net
```

Incoming messages are JSON objects:
```json
{
  "event": "message",
  "device_id": "628123456789@s.whatsapp.net",
  "payload": {
    "from": "628987654321@s.whatsapp.net",
    "message": "Hello!",
    "timestamp": "2026-05-07T12:00:00Z"
  }
}
```

**Available events:** `message`, `message.ack`, `message.reaction`, `message.revoked`, `message.edited`, `group.participants`, `group.joined`, `newsletter.joined`, `call.offer`

---

## 10. Apps Script Integration Guide

This section covers how to call the GoWA API from the Google Apps Script backend.

### Workflow: Connect a WhatsApp Account for a User

```javascript
const SERVER = "https://graphicy-wa--testhimrealty.replit.app";

// Step 1 — Check server is alive
function checkServer() {
  const res = UrlFetchApp.fetch(SERVER + "/health");
  return res.getResponseCode() === 200;
}

// Step 2 — Add a device slot for the user
function addDevice(deviceName) {
  const res = UrlFetchApp.fetch(SERVER + "/devices", {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({ name: deviceName }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
  // Returns: { results: { id: "...", ... } }
}

// Step 3 — Get QR login URL to show user
function getLoginUrl(deviceId) {
  return SERVER + "/devices/" + encodeURIComponent(deviceId) + "/login";
  // Open this URL in a browser to show QR code
}

// Step 4 — OR use pairing code (no browser needed)
function getPairingCode(deviceId, phone) {
  const res = UrlFetchApp.fetch(
    SERVER + "/devices/" + encodeURIComponent(deviceId) + "/login/code",
    {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ phone: phone }),
      muteHttpExceptions: true
    }
  );
  return JSON.parse(res.getContentText());
  // Returns: { results: { code: "ABCD-EFGH" } }
}

// Step 5 — Check device is connected
function getDeviceStatus(deviceId) {
  const res = UrlFetchApp.fetch(
    SERVER + "/devices/" + encodeURIComponent(deviceId) + "/status",
    { muteHttpExceptions: true }
  );
  return JSON.parse(res.getContentText());
  // Returns: { results: { is_connected: true, is_logged_in: true } }
}
```

---

### Sending Messages from Apps Script

Always include `X-Device-Id` header:

```javascript
// Send a plain text message
function sendTextMessage(deviceId, phone, text) {
  const res = UrlFetchApp.fetch(SERVER + "/send/message", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({ phone: phone, message: text }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}

// Send an image by URL
function sendImage(deviceId, phone, imageUrl, caption) {
  const res = UrlFetchApp.fetch(SERVER + "/send/image", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({ phone: phone, url: imageUrl, caption: caption }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}

// Send a video by URL
function sendVideo(deviceId, phone, videoUrl, caption) {
  const res = UrlFetchApp.fetch(SERVER + "/send/video", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({ phone: phone, url: videoUrl, caption: caption }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}

// Send audio by URL
function sendAudio(deviceId, phone, audioUrl) {
  const res = UrlFetchApp.fetch(SERVER + "/send/audio", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({ phone: phone, url: audioUrl }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}

// Send location
function sendLocation(deviceId, phone, lat, lng, address) {
  const res = UrlFetchApp.fetch(SERVER + "/send/location", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({
      phone: phone,
      latitude: lat,
      longitude: lng,
      address: address
    }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}

// Send contact card
function sendContact(deviceId, phone, contactName, contactPhone) {
  const res = UrlFetchApp.fetch(SERVER + "/send/contact", {
    method: "post",
    contentType: "application/json",
    headers: { "X-Device-Id": deviceId },
    payload: JSON.stringify({
      phone: phone,
      contact_name: contactName,
      contact_phone: contactPhone
    }),
    muteHttpExceptions: true
  });
  return JSON.parse(res.getContentText());
}
```

---

### Blast Messaging — Recommended Pattern

Your Apps Script already has rate limiting built in. Here is the recommended pattern for a blast campaign:

```javascript
function runBlast(deviceId, contacts, template, settings) {
  let sent = 0, failed = 0;
  const minGap = settings.minGap * 1000; // ms
  const maxGap = settings.maxGap * 1000;

  for (const contact of contacts) {
    if (contact.status === "sent") continue;

    const msg = template.body.replace("{{name}}", contact.name);
    const result = sendTextMessage(deviceId, contact.phone, msg);

    const ok = result.code === "SUCCESS";
    contact.status = ok ? "sent" : "failed";

    if (ok) sent++; else failed++;

    // Random delay between messages to avoid spam detection
    const delay = minGap + Math.random() * (maxGap - minGap);
    Utilities.sleep(delay);
  }

  return { sent, failed };
}
```

---

### AI Auto-Reply Pattern

Use the Groq AI endpoint to generate smart replies to incoming messages:

```javascript
function generateAutoReply(incomingMessage, context) {
  const res = UrlFetchApp.fetch(SERVER + "/api/ai/chat", {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({
      message: incomingMessage,
      system_prompt: context || "You are a helpful WhatsApp assistant. Reply concisely in the same language as the user.",
      conversation_history: []
    }),
    muteHttpExceptions: true
  });
  const data = JSON.parse(res.getContentText());
  return data.results?.response || "";
}
```

---

### Device Linking Flow (Multi-Tenant)

Each user in your Apps Script system gets their own `device_id`. Store it in the DeviceOwnership sheet and use it for all their API calls:

```javascript
// Complete flow: register user's WhatsApp device
function onboardUserDevice(username, phoneNumber) {
  // 1. Create device slot
  const addResult = addDevice("device_" + username);
  if (!addResult.results?.id) throw new Error("Failed to add device");
  const deviceId = addResult.results.id;

  // 2. Get pairing code
  const codeResult = getPairingCode(deviceId, phoneNumber);
  const pairingCode = codeResult.results?.code;

  // 3. Save device to Apps Script sheet
  linkDeviceAction({ username: username, deviceId: deviceId });

  return {
    deviceId: deviceId,
    pairingCode: pairingCode,
    instruction: "Enter this code in WhatsApp → Linked Devices → Link with phone number"
  };
}
```

---

### Webhook Setup

To receive incoming messages from your users' WhatsApp accounts, configure a webhook URL when starting the server:

```
WHATSAPP_WEBHOOK=https://script.google.com/macros/s/<your-script-id>/exec
```

The server will POST to this URL with payloads like:
```json
{
  "event": "message",
  "device_id": "628123456789@s.whatsapp.net",
  "payload": {
    "from": "628987654321@s.whatsapp.net",
    "message": "Hi, I need help",
    "timestamp": "2026-05-07T12:00:00Z"
  }
}
```

In your Apps Script `doPost`, read `e.postData.contents` and route by `event` and `device_id`.

---

## 11. Common Response Format

All API responses follow this structure:

```json
{
  "status": 200,
  "code": "SUCCESS",
  "message": "Human readable message",
  "results": { }
}
```

**Error codes:**

| HTTP | `code` | Meaning |
|---|---|---|
| 400 | `BAD_REQUEST` | Invalid input |
| 401 | `UNAUTHORIZED` | Wrong basic-auth credentials |
| 404 | `NOT_FOUND` | Device or resource not found |
| 422 | `VALIDATION_ERROR` | Request validation failed |
| 500 | `INTERNAL_SERVER_ERROR` | Server-side error |
| 503 | `SERVICE_UNAVAILABLE` | WhatsApp not connected |

---

## Quick Reference — Phone Number Format

- Always use **international format without `+`**
- Examples: `628123456789` (Indonesia), `14155552671` (US), `447700900123` (UK)
- For groups: `120363xxxxx@g.us`
- For users: `628123456789@s.whatsapp.net`

When sending via API, pass just the number part (e.g. `628123456789`) — the server appends `@s.whatsapp.net` automatically.
