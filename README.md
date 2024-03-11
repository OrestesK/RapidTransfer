# Rapid Transfer System

Rapid Transfer System is a simple system for transferring files between users using Go and a PostgreSQL database.

## Future Features

- Auto ziping folders
- Imroving file pathing
- ~~changing arg parsing to accept no defaults~~
- ~~Improvement of SQL usage, increase efficiency~~
- ~~move away from mac addresses and into password authentication~~
- ~~hash information~~
- able to send files from anywhere
- choose where files are recieved
- ~~more documentation on how it works~~
- ~~better error handling~~

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

   ```
   git clone https://github.com/your-username/rapid-transfer-system.git

   Navigate to the project directory:

   Download binary using make
   make [windows, linux, darwin] 
   ```
2. Connect your SQL credentials inside of the init_database file (Will later change to make it easier to use)

Usage
Command Line Flags
```
	-send user # Used to send file to user, must use -file path flag to specify the file
	-file path # Used to specify the path to the file you are sending, must be used with -send
	-add user_id # Used to add a friend, user_id is the id you retrieve when you use -info
	-inbox # Used to retrive information about files you have yet to accept
	-delete filename # Used to remove a file from your inbox
	-boot friend_id # Used to remove a friend from your friends list
	-recieve file # Used to accept a file being sent to you
	-friends # Used to list all of your friends and their friend id
	-info # Used to display your account information
```
Examples (After building the binary)
Send a File:

Rapid -send adam -file memes.png

Receive a File:

Rapid -recieve memes.png

View Pending Transfers:

Rapid -inbox

Delete a Friend:

Rapid -add adam

Add a Friend:

Rapid -delete adam

This project is licensed under the MIT License - see the LICENSE file for details.
