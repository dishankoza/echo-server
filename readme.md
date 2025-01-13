# Async Echo Server  

A showcase of different server architectures and the evolution of I/O handling techniques. This repository demonstrates implementations of synchronous, asynchronous, and event-driven TCP servers, culminating in an efficient implementation using `epoll` and `io_uring`.

---

## 🚀 Overview  

This repository contains:  
1. **Synchronous TCP Server:** A basic server that handles one client at a time.  
2. **Asynchronous TCP Server with Goroutines:** Utilizes goroutines to handle multiple connections but can become resource-intensive at scale.  
3. **Event-Driven TCP Server with Epoll:** A single-threaded architecture leveraging `epoll` for I/O multiplexing to handle thousands of concurrent connections efficiently.  

---

## 🛠️ Technologies  

- **Go**: Programming language for server implementation.  
- **syscall**: Direct system calls to interact with the OS (used for `epoll`).  
- **io_uring**: Cutting-edge Linux I/O interface for reduced system call overhead (in progress).  

---

## 📂 Repository Structure  

- `sync/`  
  - Contains the synchronous server implementation.  
- `async-goroutines/`  
  - Implements the server with Goroutines.  
- `async-epoll/`  
  - Implements the server with `epoll` and event-driven architecture.  
- `core/`  
  - Utilities for reading and handling client data (RESP protocol).  

---

## 🧠 Concepts  

1. **Synchronous Blocking I/O**  
   Each client blocks the server until the operation completes, limiting scalability.  

2. **Multi-threading with Goroutines**  
   Each connection spawns a new goroutine, improving concurrency but consuming more memory as clients increase.  

3. **Single-Threaded Event-Driven I/O (Epoll)**  
   A single thread manages thousands of connections using an event-driven approach, waiting for I/O readiness signals.  

4. **io_uring**  
   The future of I/O handling with minimal overhead, allowing direct interaction with kernel queues.  

---

## 📝 How It Works  

### Event Loop with Epoll  
- A single thread monitors multiple file descriptors (connections) using `epoll`.  
- When a file descriptor signals readiness (e.g., data received), the server processes it.  
- This avoids blocking and reduces the overhead of thread management.  

---

## 🚀 Run the Server  

### Prerequisites  
- Go installed on your system.  
- A Linux-based system (required for `epoll` and `io_uring`).  

### Steps  
1. Clone the repository:  
   ```bash
   git clone https://github.com/dishankoza/echo-server.git
   cd echo-server
