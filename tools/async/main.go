package async

// https://github.com/rafaeldias/async/blob/master/async.go

// import (
// 	"context"
// )
//
// func executor(ctx context.Context, errCh chan error, fn func() error) {
// }
//
// func Semaphore(fns func() error...) error {
// 	errCh := make(chan error, 0, len(fns))
// 	ctx, cancel := context.WithCancel(context.Background())
//
// 	defer cancel()
//
// 	for _, fn := range fns {
// 		go func() {
// 			if err := fn(); err != nil {
// 				errCh <- err
// 				return
// 			}
// 		}()
// 	}
//
// 	select{
// 	case err := <- errCh:
// 		return error
// 	}
// }
//
//
//
// // sync
// for _, fn := range fns {
// 	if err := fn(); err != nil {
// 		errCh <- err
// 		return
// 	}
// }
//
//
//
//
// // async
// wg.Add(len(fns))
// for _, fn := range fns {
// 	go func() {
// 		if err := fn(); err != nil {
// 			errCh <- err
// 			return
// 		}
// 	}()
// }
