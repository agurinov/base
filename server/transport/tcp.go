package transport

import (
	"io"
	"net"
	"time"
	// "github.com/boomfunc/log"
)

var (
	// TODO parametrize
	timeout = time.Second * 5
)

type tcp struct {
	listener *net.TCPListener
	inputCh  chan io.ReadWriteCloser
	errCh    chan error
}

func (tr *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

// https://habr.com/company/mailru/blog/331784/
// before 3.3.1
func (tr *tcp) Serve() {
	for {
		conn, err := tr.listener.AcceptTCP()
		if err != nil {
			// handle error
			tr.errCh <- err
			continue
		}

		// // handle successful connection
		// // TODO maybe send connections only when the caller starts to write to it?
		// // TODO maybe send connections also when worker can be fetched?
		// go func(conn *net.TCPConn) {
		//
		// 	raw, _ := conn.SyscallConn()
		//
		// 	log.Debug("PSEUDO CHECK - READY", raw.Read(tcpDetectRead))
		//
		//
		// }(conn)

		conn.SetDeadline(time.Now().Add(timeout))
		tr.inputCh <- conn
	}
}





// package transport
//
// import (
// 	"io"
// 	"net"
// 	"time"
//
// 	"github.com/boomfunc/log"
// 	// TODO deal with it!!!!!!
// 	"github.com/mailru/easygo/netpoll"
// )
//
// var (
// 	// TODO parametrize
// 	timeout = time.Second * 5
// )
//
// type tcp struct {
// 	listener *net.TCPListener
// 	inputCh  chan io.ReadWriteCloser
// 	errCh    chan error
// }
//
// func (tr *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
// 	tr.inputCh = inputCh
// 	tr.errCh = errCh
// }
//
// // https://habr.com/company/mailru/blog/331784/
// // before 3.3.1
// func (tr *tcp) Serve() {
// 	poller, err := netpoll.New(nil)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Debug("LISTENING")
// 	for {
// 		conn, err := tr.listener.AcceptTCP()
// 		if err != nil {
// 			// handle error
// 			tr.errCh <- err
// 			continue
// 		}
//
// 		log.Debug("CONN ARRIVED")
//
// 		// Get netpoll descriptor with EventRead|EventEdgeTriggered.
// 		desc := netpoll.Must(netpoll.HandleRead(conn))
// 		poller.Start(desc, func(ev netpoll.Event) {
// 			// We spawn goroutine here to prevent poller wait loop
// 			// to become locked during receiving packet from ch.
// 			log.Debug("CONN READY", conn)
// 			go func() {
// 				conn.SetDeadline(time.Now().Add(timeout))
// 				tr.inputCh <- conn
// 			}()
//
// 			// TODO poller stop
// 		})
//
// 		log.Debug("GO NEXT")
//
// 		// handle successful connection
// 		// TODO maybe send connections only when the caller starts to write to it?
// 		// TODO maybe send connections also when worker can be fetched?
// 		// go func(conn *net.TCPConn) {
// 		//
// 		// 	raw, _ := conn.SyscallConn()
// 		//
// 		// 	raw.Read(tcpOnRead)
// 		//
// 		// 	// conn.SetDeadline(time.Now().Add(timeout))
// 		// 	tr.inputCh <- conn
// 		//
// 		// }(conn)
//
// 	}
// }
