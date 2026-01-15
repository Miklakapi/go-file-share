# Go-File-Share

![license](https://img.shields.io/badge/license-MIT-blue)
![linux](https://img.shields.io/badge/os-Linux-green)
![language](https://img.shields.io/badge/language-Go_1.25.1-blue)
![version](https://img.shields.io/badge/version-1.0.0-success)
![status](https://img.shields.io/badge/status-production-green)

An overengineered real-time file sharing server written in Go.  
The project allows temporary file exchange between devices using password-protected rooms, direct streaming, and Server-Sent Events (SSE).

This project was created mainly as a **playground for architecture, system design, and low-level implementation details** rather than as a minimal solution.

## Table of Contents

-   [General info](#general-info)
-   [Architecture](#architecture)
-   [Technologies](#technologies)
-   [Setup](#setup)
-   [Features](#features)
-   [Status](#status)

## General info

Go-File-Share is a temporary file-sharing server that allows users to exchange files between devices without persistent storage.

The system is based on **rooms**:

-   each room is protected by a password,
-   users can upload and download files within a room,
-   rooms have a limited lifetime and are automatically removed after expiration.

In addition to room-based sharing, the application supports **direct file transfer between users**:

-   one user generates a temporary connection code,
-   another user sends a file using that code,
-   the file is streamed through the server **without being saved to disk**.

The frontend is written in **plain JavaScript**, without any frameworks, and communicates with the backend using HTTP, SSE, and streaming endpoints.

> **Important:**  
> All files and database data are **intentionally wiped on every application startup**.  
> This is a conscious design decision. Each restart returns the system to a **clean, zero-state**.

<p align="center" width="100%">
    <img width="100%" src="https://github.com/Miklakapi/go-file-share/blob/master/README_IMAGES/room_list.png"> 
    <img width="100%" src="https://github.com/Miklakapi/go-file-share/blob/master/README_IMAGES/room.png"> 
    <img width="50%" src="https://github.com/Miklakapi/go-file-share/blob/master/README_IMAGES/direct_transfer.png">
</p>

## Architecture

This project intentionally uses an **overengineered, hexagonal (ports & adapters) architecture**.

Core characteristics:

-   strict separation of domain logic and infrastructure,
-   heavy use of interfaces and ports,
-   multiple interchangeable adapters.

### Adapters

Currently implemented:

-   **In-memory repository** – fast, volatile storage
-   **SQLite repository** – persistent storage (reset on startup)
-   **Redis repository** – distributed, ephemeral storage adapter

### Internal systems

The application also includes several custom-built components:

-   custom **event bus** used to propagate domain events,
-   **SSE layer** for real-time room list updates,
-   background **scheduler/job** that periodically removes expired rooms,
-   custom **database migration tool** written in plain Go,
-   file streaming between users using `io.Reader` / `io.Writer` without buffering files on disk,
-   simple file-based logging.

The main goal was to avoid external abstractions and frameworks wherever possible and rely on **plain Go** and **standard library primitives**.

## Technologies

Project is created with:

-   Go 1.25.1
-   Gin 1.11.0 – HTTP server and routing
-   SQLite – local database
-   Redis – in-memory data store
-   Server-Sent Events (SSE) – real-time updates
-   Plain JavaScript – frontend (no frameworks)

## Setup

### Server

1. Install dependencies:
    ```
    go mod tidy
    ```
2. Run the application:
    ```
    go run ./cmd/sqlite-app/main.go
    ```
3. Server starts on:
    ```
    http://localhost:8080
    ```

### Frontend

The frontend is served directly by the Go server.

Open your browser and navigate to:

```
http://localhost:8080
```

## Features

-   Temporary file sharing between devices
-   Password-protected rooms
-   Room expiration with automatic cleanup
-   Real-time room list updates via SSE
-   Direct file transfer between users using connection codes
-   Streaming file transfer without saving files on the server
-   Hexagonal architecture (ports & adapters)
-   Multiple repository implementations (RAM, SQLite, Redis)
-   Custom event bus
-   Custom database migration tool
-   Plain Go backend & plain JavaScript frontend

## Status

The project's development has been completed.
