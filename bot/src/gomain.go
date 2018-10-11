package main

import (
	"fmt"
	"math/rand"
	"strings"
	//"sync"
	//"sync/atomic"
	"time"
)

const benchmark = false
const N_USERS int = 20
const MAX_MESSAGES int = 256

var n_messages = 0

type MESSAGE struct {
	msg  string
	src  string
	dst  string
	next *MESSAGE
}

type CHATBOT struct {
	n_messages   int
	messages     [MAX_MESSAGES]MESSAGE
	last_message *MESSAGE
	free_message *MESSAGE
	mutex        chan int
}

type USER struct {
	name      string
	is_online int
	next      *USER
}

type CHANNEL struct {
	n_users   int
	users     [N_USERS]USER
	last_user *USER
	free_user *USER
	mutex     chan int
}

func chatbot_start(bot *CHATBOT) {
	for i := 0; i < MAX_MESSAGES-1; i++ {
		bot.messages[i].next = &bot.messages[i+1]
	}
	bot.free_message = &bot.messages[0]
	bot.mutex = make(chan int, 1)
	bot.mutex <- 0
}

func chatbot_login(bot *CHATBOT, username string) {
	<-bot.mutex

	if !benchmark {
		fmt.Printf("%s chatbot login\n", username)
	}
	var prev *MESSAGE = nil
	for node := bot.last_message; node != nil; {
		if strings.Compare(node.dst, username) == 0 {
			if !benchmark {
				fmt.Printf("Hello %s, %s left this message: \"%s\"\n", node.dst, node.src, node.msg)
			}
			n_messages++
			next := node.next
			if prev != nil {
				prev.next = node.next
			}
			if node == bot.last_message {
				bot.last_message = node.next
			}

			node.next = bot.free_message
			bot.free_message = node

			bot.n_messages--
			node = next
		} else {
			prev = node
			node = node.next
		}
	}
	bot.mutex <- 0
}

func chatbot_leave_message(bot *CHATBOT, msg string, src string, dst string) {

	if dst == "" || msg == "" {
		return
	}

	if !benchmark {
		fmt.Printf("%s leaving message for %s\n", src, dst)
	}

	<-bot.mutex

	if bot.free_message != nil {
		node := bot.free_message
		bot.free_message = bot.free_message.next
		node.next = bot.last_message
		bot.last_message = node

		bot.n_messages++

		node.msg = msg
		node.src = src
		node.dst = dst
	} else {
		if !benchmark {
			fmt.Printf("Unable to store message: Out of memory.")
		}
	}

	bot.mutex <- 0
}

func channel_start(channel *CHANNEL) {
	for i := 0; i < N_USERS-1; i++ {
		channel.users[i].next = &channel.users[i+1]
	}
	channel.free_user = &channel.users[0]
	channel.mutex = make(chan int, 1)
	channel.mutex <- 0
}

func channel_login(channel *CHANNEL, name string) {

	<-channel.mutex
	if channel.free_user != nil {
		node := channel.free_user
		channel.free_user = channel.free_user.next
		node.next = channel.last_user
		channel.last_user = node

		channel.n_users++

		node.name = name
		node.is_online = 1
	} else {
		if !benchmark {
			fmt.Printf("Channel is full.\n")
		}
	}
	channel.mutex <- 0
}

func channel_logout(channel *CHANNEL, name string) {

	<-channel.mutex
	var prev *USER = nil
	for node := channel.last_user; node != nil; node = node.next {
		if strings.Compare(node.name, name) == 0 {
			if !benchmark {
				fmt.Printf("%s logged out\n", name)
			}
			node.is_online = 0
			if prev != nil {
				prev.next = node.next
			}
			if node == channel.last_user {
				channel.last_user = node.next
			}

			node.next = channel.free_user
			channel.free_user = node

			channel.n_users--
			break
		}

		prev = node
	}
	channel.mutex <- 0
}

func channel_get_username(channel *CHANNEL) string {

	name := ""
	<-channel.mutex

	r := 0
	if channel.n_users > 0 {
		r = rand.Int() % channel.n_users
	}

	for node := channel.last_user; (node != nil) && (1+r > 0); node = node.next {
		name = node.name
		r -= 1
	}

	channel.mutex <- 0
	return name
}

var chatbot CHATBOT
var channel CHANNEL

func user_thread(data int) {
	//rand.Seed(time.Now().UTC().UnixNano())

	name := ""
	dst := ""
	for i := 0; i < 3; i++ {
		name += string(97 + (rand.Int() % 26))
	}
	message := "Hello from " + name + "!"

	dst = message
	message = dst
	for {
		channel_login(&channel, name)
		chatbot_login(&chatbot, name)
		time.Sleep(time.Millisecond * time.Duration(rand.Int()%200))

		chatbot_leave_message(&chatbot, message, name, dst)
		channel_logout(&channel, name)

		dst = channel_get_username(&channel)
		time.Sleep(time.Millisecond * time.Duration(rand.Int()%200))
	}
}

func main() {
	chatbot_start(&chatbot)
	channel_start(&channel)

	for i := 0; i < N_USERS; i++ {
		go user_thread(i)
	}

	go func() {
		time.Sleep(time.Millisecond * 10 * 1000)
		//fmt.Printf("messages sent: %d\n", n_messages)
	}()

	for {

	}

}
