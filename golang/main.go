type Ttype struct {
	ID         int
	CreateTime time.Time
	FinishTime time.Time
	Result     []byte
}

func main() {
	const taskCount = 10

	taskCreturer := func(c chan<- Ttype) {
		for {
			t := Ttype{
				ID:         int(time.Now().Unix()),
				CreateTime: time.Now(),
			}

			if time.Now().Nanosecond()%2 > 0 {
				t.Result = []byte("Some error occured")
			} else {
				t.Result = []byte("")
			}
			c <- t
		}
	}

	taskWorker := func(t Ttype) Ttype {
		if time.Since(t.CreateTime) < 20*time.Second {
			t.Result = []byte("Task has been successed")
		} else {
			t.Result = []byte("Something went wrong")
		}
		t.FinishTime = time.Now()
		time.Sleep(150 * time.Millisecond)

		return t
	}

	taskSorter := func(t Ttype) (Ttype, error) {
		if string(t.Result) == "Task has been successed" {
			return t, nil
		}
		return Ttype{}, fmt.Errorf("Task id %d create time %s, error %s", t.ID, t.CreateTime, t.Result)
	}

	taskChan := make(chan Ttype, taskCount)
	doneTasks := make(chan Ttype, taskCount)
	undoneTasks := make(chan error, taskCount)

	go taskCreturer(taskChan)

	for i := 0; i < taskCount; i++ {
		go func() {
			for t := range taskChan {
				t = taskWorker(t)
				sorted, err := taskSorter(t)
				if err != nil {
					undoneTasks <- err
				} else {
					doneTasks <- sorted
				}
			}
		}()
	}

	results := make(map[int]Ttype)
	errors := []error{}

	for i := 0; i < taskCount; i++ {
		select {
		case t := <-doneTasks:
			results[t.ID] = t
		case e := <-undoneTasks:
			errors = append(errors, e)
		}
	}

	fmt.Printf("Successful tasks: %d\n", len(results))
	for _, t := range results {
		fmt.Printf("Task %d: create time %s, finish time %s\n", t.ID, t.CreateTime, t.FinishTime)
	}

	fmt.Printf("Failed tasks: %d\n", len(errors))
	for _, e := range
