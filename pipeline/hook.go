package pipeline

// type Hook func() error
// type Hooks []Hook
//
// func (hks *Hooks) sync(silent bool) []error {
// 	// exec variants:
// 	// 1. Run each by order, stops when err raised
// 	// 2. Run each by order, collect all errors
// 	var errors []error
//
// 	for i, hook := hks {
// 		if err := hook(); err != nil {
// 			errors = append(errors, err)
//
// 			if !silent {
// 				return errors
// 			}
// 		}
// 	}
//
// 	return errors
// }
//
// func (hks *Hooks) async(silent bool) {
// 	// exec variants:
// 	// 3. Run all async, collect all errors
// 	// 4. Run all async, stops when err raised (stop with context)
// }
