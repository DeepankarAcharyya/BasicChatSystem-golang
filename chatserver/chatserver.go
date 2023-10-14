package chatserver

import (
	"log"
	"math/rand"
	sync "sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

// messageHandle will hold slice of messageUnits
type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex // add a mutex variable to handle async read write operations
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

// define ChatService
func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {

	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// receive messages
	go receiveFromStream(csi, clientUniqueCode, errch)

	// send messages
	go sendToStream(csi, clientUniqueCode, errch)

	// On recceiving error through channel errch, the connection will be terminated
	return <-errch
}

// receive messages
func receiveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	// infinite loop to receive the messages
	for {
		mssg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from the client :: %v", err)
			errch_ <- err
		} else {
			messageHandleObject.mu.Lock()
			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        mssg.Name,
				MessageBody:       mssg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode_,
			})
			messageHandleObject.mu.Unlock()
			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}
	}
}

// send messages
func sendToStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	// implement a loop
	for {
		//loop through messages in MQue
		for {
			time.Sleep(599 * time.Millisecond)
			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObject.MQue[0].ClientName
			message4Client := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			//send message to designated client
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(
					&FromServer{
						Name: senderName4Client,
						Body: message4Client,
					},
				)

				if err != nil {
					errch_ <- err
				}

				messageHandleObject.mu.Lock()
				if len(messageHandleObject.MQue) > 1 {
					// deleting the sent message from the queue
					messageHandleObject.MQue = messageHandleObject.MQue[1:]
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}
				messageHandleObject.mu.Unlock()
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
