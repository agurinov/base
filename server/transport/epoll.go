// +build linux

package transport

// import (
// 	"syscall"
// )
//
//
// func NewPoller() Poller {
//
// }
//
// type epoll struct {
// 	epfd int
// }
//
// func (p *epoll) Add(fd uintptr) error {
// 	var event syscall.EpollEvent
// 	event.Events = syscall.EPOLLIN | syscall.EPOLLET
// 	event.Fd = int32(fd)
//
// 	return syscall.EpollCtl(p.epfd, syscall.EPOLL_CTL_ADD, int(fd), &event)
// }
//
// func (p *epoll) Del(fd uintptr) error {
// 	return syscall.EpollCtl(p.epfd, syscall.EPOLL_CTL_DEL, int(fd), nil)
// }
//
// func (p *epoll) Wait(fd uintptr) error {
// 	for {
// 		nevents, err := syscall.EpollWait(p.epfd, events[:], -1)
// 		log.Debug("EPOLL TRIGGERED", nevents)
// 		if err != nil {
// 			return err
// 		}
//
// 		// TODO
//
// 	}
// }
