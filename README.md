# Rapid Transfer System

Rapid Transfer System is a simple system for transferring files between users using Go and a PostgreSQL database.

## Future Features

- Auto ziping folders
- Imroving file pathing
- changing arg parsing to accept no defaults
- Improvement of SQL usage, increase efficiency 
- move away from mac addresses and into password authentication
- hash information
- able to send files from anywhere
- choose where files are recieved
- more documentation on how it works
- better error handling

## Overview

The Rapid Transfer System is designed for easy and fast file transfers between users. It utilizes Go for the backend logic, a PostgreSQL database for user information, and a network package for handling file transfers.

## Features

- Initialize and manage the underlying PostgreSQL database.
- Handle account startup tasks.
- Manage user details and friends.
- Send and receive files between users.
- View and manage pending transfers.

## Getting Started

### Prerequisites

Make sure you have Go installed on your system. You can download it [here](https://golang.org/dl/).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/rapid-transfer-system.git

    Navigate to the project directory:

    bash
   ```

cd rapid-transfer-system

Build the project:

bash

    go build

Usage
Command Line Flags

    -s or --send: Specify the file to send.
    -p or --path: Specify the path to the file to send.
    --friend: Specify the friend's username for sending or friend-related operations.
    -r or --receive: Specify the file to receive.
    -d or --delete: Specify the friend to delete.
    --pend: View pending transfers.

Examples
Send a File:

./rapid-transfer-system -s filename.txt --friend friend_username

Receive a File:

./rapid-transfer-system -r filename.txt

View Pending Transfers:

./rapid-transfer-system --pend

Delete a Friend:

./rapid-transfer-system -d friend_username

This project is licensed under the MIT License - see the LICENSE file for details.

sql

Replace `your-username` in the clone URL with your actual GitHub username, and add any specific details, usage instructions, or guidelines that might be relevant to your project.
