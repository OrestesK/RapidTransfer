# Rapid
Originally this project was called RapidTransfer and was created within 24 hours. It was capable of adding friends, and sharing just files to eachother using peer to peer networking. I took that project and changed/added some essential features
- Better documentation, pathing and allowing users to send folders
- cleaner queries
- Cloud is used instead of p2p to save files
- Encyption based on AES and hashing of all personal information
- More functionality

## Why Cloud?
Transfering files between users using p2p is cool but there are some issues that come with that. Here are some of the main ones that I had thought about
1. Disconnection while sending files due to exiting the program or turning off device
2. Need for both users to be active at the same time
3. Higher level of understanding needed to be able to reproduce the code

Implementing Cloud allows these added benifits
- Allow for easy access of files
- Both users do not have to be online
- Once the sender sends the file, they can remove it locally from their system
- Easier to read the code (Cloud made things very easy to write)

## Now
Rapid is a faster and better way to share files between friends without having to worry about other people getting access to them. Each file is encrypted and uploaded to the cloud with a unique index. User information is stored using SQL aswell as information about the transaction that allows the files to be decrypted once downloaded back onto a users machine. User authentication exists which makes sure that nobody can gain access to your account. Mac addresses are used as a double authentication to make sure that a user would need to be on the same device that the account was created on to access their account.

## Features
- File sharing between users
- Friend functionality so users can add and remove friends
- Built in encryption and decryption using AES with randomized keys in order to ensure uniqueness
- Able to view ones own inbox and choose to remove or accept incoming files

## Getting Started
1. Make sure you have Go installed on your system. You can download it [here](https://golang.org/dl/).
2. Clone the repository ```git clone https://github.com/your-username/rapid-transfer-system.git```
3. Download the **CORRECT** binary of MEGACMD and make sure it is **included** within your path [here](https://github.com/t3rm1n4l/megacmd/releases/tag/0.016)
4. Create a file ~/.megacmd.json inside of your user folder. Here is an example **The user information will either be provided by the hoster or you will have to create yourself**
```
{
    "User" : "MEGA_USERNAME",
    "Password" : "MEGA_PASSWORD",
    "DownloadWorkers" : 4,
    "UploadWorkers" : 4,
    "SkipSameSize" : true,
    "Verbose" : 1
}
```
5. Enter in the SQL credentials inside of the private.go file **This will either be provided by the hoster or you will have to create yourself**

## Usage
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
Examples **After completing all the steps to get the binary working**
```
Send a File:
Rapid -send adam -file memes.png
Receive a File:
Rapid -recieve memes.png
View Pending Transfers:
Rapid -inbox
Delete a Friend:
Rapid -add adam
```

Add a Friend:

Rapid -delete adam

This project is licensed under the MIT License - see the LICENSE file for details.
