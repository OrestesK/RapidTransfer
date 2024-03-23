package main

import (
	"Rapid/src/cloud"
	"Rapid/src/database"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

// displays friends list
func displayFriends(user int) {
	friendsList := database.GetFriendsList(user)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Friend Code"})
	for _, friend := range friendsList {
		fmt.Println(friend.Name, friend.FriendCode)
		t.AppendRows([]table.Row{
			{friend.Name, friend.FriendCode},
		})
	}
	t.Render()
}

// displays inbox
func displayInbox(user int) {
	inbox := database.GetPendingTransfers(user)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"From", "File Name"})
	for _, transaction := range inbox {
		t.AppendRows([]table.Row{
			{transaction.From_user, transaction.File_name},
		})
	}
	t.Render()
}

func appStartup() {
	var user int = 1
	result, _ := os.UserHomeDir()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login [Username] [Password] {Login to a users account}",
				Action: func(c *cli.Context) error {
					err := database.SetCurrentUsersId(c.Args().First(), c.Args().Get(1))
					if err != nil {
						return err
					}
					user = database.GetCurrentId()
					fmt.Printf("Currently Logged in as %s\n", database.GetUserNameByID(user))
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create [Username] [Password] {Create a users account}",
				Action: func(c *cli.Context) error {
					err := database.CreateAccount(c.Args().First(), c.Args().Get(1))
					if err != nil {
						return err
					}

					err = database.SetCurrentUsersId(c.Args().First(), c.Args().Get(1))
					if err != nil {
						return err
					}

					user = database.GetCurrentId()
					fmt.Printf("Currently Logged in as %s\n", database.GetUserNameByID(user))

					return nil
				},
			},
			{
				Name:  "help",
				Usage: "help {Displays all commands and information}",
				Action: func(c *cli.Context) error {
					cli.ShowAppHelp(c)
					return nil
				},
			},
			{
				Name:    "user",
				Usage:   "user, u {Displays user information}",
				Aliases: []string{"u"},
				Action: func(c *cli.Context) error {
					fmt.Printf("| Username   %s | Friend code   %s |\n", database.GetUserNameByID(user), database.GetUserFriendCode(user))
					return nil
				},
			},
			{
				Name:    "send",
				Usage:   "send, s [User] [Filepath] {Will send user file/folder}",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					err, result := cloud.UploadToMega(c.Args().Get(1), user, c.Args().First())
					if err != nil {
						fmt.Println(err)
					}
					if result {
						fmt.Println("File has been sent and will be waiting to be accepted")
					} else {
						fmt.Println("The requested user either does not exist or is not added")
					}
					return nil
				},
			},
			{
				Name:  "inbox",
				Usage: "inbox recieve, r [Filename] | inbox remove, rm [Filename] | inbox list, l {Handles inbox functionality}",
				Subcommands: []*cli.Command{
					{
						Name:    "recieve",
						Aliases: []string{"r"},
						Usage:   "inbox recieve, r [Filename] {Recieves file from inbox}",
						Action: func(c *cli.Context) error {
							fmt.Println("Key is:", c.String("key"))
							_, result := cloud.DownloadFromMega(user, c.Args().First(), "")
							if result {
								fmt.Println("File has been received")
							} else {
								fmt.Println("Filename or item does not not exist within your inbox")
							}
							return nil
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "inbox remove, rm [Filename] {Removes file from inbox}",
						Action: func(c *cli.Context) error {
							result, err := cloud.DeleteFromMega(user, c.Args().First())
							if err != nil {
								fmt.Println(err)
							}
							if result {
								fmt.Println("File has been deleted")
							} else {
								fmt.Println("Could not delete the file from the inbox")
							}
							return nil
						},
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "inbox list, l {Lists all messages in inbox}",
						Action: func(c *cli.Context) error {
							displayInbox(user)
							fmt.Println("Inbox has been displayed")
							return nil
						},
					},
				},
			},
			{
				Name:  "friend",
				Usage: "friend add, a [Friend id] | friend remove, rm [Username] | friend list, l {Handles friend functionality}",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "friend add, a [Friend id] {Adds friend}",
						Action: func(c *cli.Context) error {
							result, err := database.AddFriend(c.Args().First(), user)
							if err != nil {
								fmt.Println(err)
							}

							if !result {
								fmt.Println("There does not exist a user with that friend code.")

							}

							if result {
								fmt.Println("Friend has been added")
							}
							return nil
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "friend remove, rm [Username] {Removes friend}",
						Action: func(c *cli.Context) error {
							result, err := database.DeleteFriend(user, c.Args().First())
							if err != nil {
								fmt.Println("Failed to remove friend", err)
							}
							if result {
								fmt.Println("Friend has been deleted")
							}
							return nil
						},
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "friend list, l {Lists all friends}",
						Action: func(c *cli.Context) error {
							displayFriends(user)
							fmt.Println("Friends list has been displayed")
							return nil
						},
					},
				},
			},
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "key",
				Usage:       "-key, -k [Path To Key] {Specifies path to encryption key}",
				Value:       filepath.Join(result, "Rapid", "supersecretekey.txt"),
				Destination: &result,
				Aliases:     []string{"k"},
			},
		},
		EnableBashCompletion: true,
	}
	app.Suggest = true

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	appStartup()
}
