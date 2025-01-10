## Extended scheduler

The library provides the ability to create periodic jobs using the advanced schedule format used in Unix-like systems. You can create groups of buildings that will be repeated with a certain periodicity according to the format described below<br><br>

<p align="center">
    <a href="https://github.com/dimonrus/gojob/releases"><img src="https://img.shields.io/github/tag/dimonrus/gojob.svg?label=version&color=brightgreen"></a>
    <a href="https://github.com/dimonrus/gojob/actions/workflows/go.yml"><img src="https://github.com/dimonrus/gojob/actions/workflows/go.yml/badge.svg"></a>
    <a href="https://github.com/dimonrus/gojob/blob/master/LICENSE"><img src="https://img.shields.io/github/license/dimonrus/gojob"></a>
</p>

### Example of usage
```go
    // Set required period for job schedule
	// This means that all jobs in the group will be scheduled every 1 second.
	// If you don't want to call the method, the default repeat duration will be 1 minute
    gojob.SetRepeatDuration(time.Second)
    // Add to schedule you custom jobs
    gojob.Add("test.hello.job", "- * * * * * * * *", func(ctx context.Context, args ...any) error {
		panic("internal error")
		return nil
	})
    // Each job should have its own name for better understanding of the execution processes and investigation of event log entries
    gojob.Add("test.goodbye.job", "- * - - - - - - -", func(ctx context.Context, args ...any) error {
		return nil
	}, 
	// You can add your own conditions to start a job. Conditions can be combined to create complex mechanics for starting jobs
	NewCondition(OperatorAND, func() bool {
		return true
	}))
    // Once all jobs have been added, you should run the scheduler. 
	// You can add your own middlewares to provide the necessary functionality when running jobs
    gojob.Run(log.Default(), LogMiddleware, RecoverMiddleware)
```
### Schedule expression format

Schedule expression contains 9 parts. Each part responsible for special time period

```* * * * * * * * * ```<br>
```1 2 3 4 5 6 7 8 9 ```

1) Millisecond - possible values is 0-999
2) Second - possible values is 0-59
3) Minute - possible values is 0-59
4) Hour - possible values is 0-23
5) DayOfWeek - possible values is 1-7
6) DayOfMonth - possible values is 1-31
7) WeekOfMonth - possible values is 1-5
8) WeekOfYear - possible values is 1-53
9) Month - possible values is 1-12

You can specify many values using the standard syntax of Unix systems:

1) Range example: ```200-500 - - - - - - - -``` - every millisecond between 200 and 500 millisecond
2) Each of example ```- */5 * - - - - - -``` - each 5 second every minute
3) Combined example ```- - 1-30/6 - - - - - -``` - each 6th minute from 1 to 30 minutes
4) Combined groups ```- - - 1,2,3,*/6,20-23 - - - - -``` At 1,2,3, every 6th hour, and every hour from 20-23 

### Scheduler settings

1) Mode
   - GroupModeConsistently - all jobs will be scheduled consistently in order according to job sort order. Either in the order in which they are added to the scheduler
   - GroupModeAllParallel - all jobs will be scheduled in parallel mode in order according to job sort order. Either in the order in which they are added to the scheduler
   - Specify number of parallel jobs - You can specify a specific number of simultaneous jobs in a group
2) Repeat duration - determines the time interval until the next call to check conditions and start jobs
3) Middleware - define the preparatory steps before starting a jobs. These can be both logging and panic protection functions

#### If you find this project useful or want to support the author, you can send tokens to any of these wallets
- Bitcoin: bc1qgx5c3n7q26qv0tngculjz0g78u6mzavy2vg3tf
- Ethereum: 0x62812cb089E0df31347ca32A1610019537bbFe0D
- Dogecoin: DET7fbNzZftp4sGRrBehfVRoi97RiPKajV