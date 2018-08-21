package async

// import (
// 	"context"
// )
//
// func executor(ctx context.Context, errCh chan error, fn func() error) {
// }
//
// func Semaphore(fns func() error...) error {
// 	errCh := make(chan error)
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
//
// 		}()
// 	}
//
// 	select{
// 	case err := <- errCh:
// 		return error
// 	}
// }
