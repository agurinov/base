package pipeline

// req, _ := http.NewRequest("GET", "http://localhost:1111", nil)
//    go func() {
// 	   resp, err := client.Do(req)
// 	   fmt.Println("Doing http request is a hard job")
// 	   pack := struct {
// 		   r   *http.Response
// 		   err error
// 	   }{resp, err}
// 	   c <- pack
//    }()
//
//    select {
//    case <-ctx.Done():
// 	   tr.CancelRequest(req)
// 	   <-c // Wait for client.Do
// 	   fmt.Println("Cancel the context")
// 	   return ctx.Err()
//    case ok := <-c:
// 	   err := ok.err
// 	   resp := ok.r
// 	   if err != nil {
// 		   fmt.Println("Error ", err)
// 		   return err
// 	   }
//
// 	   defer resp.Body.Close()
// 	   out, _ := ioutil.ReadAll(resp.Body)
// 	   fmt.Printf("Server Response: %s\n", out)
//
//    }
//    return nil
