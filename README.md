## grpc-protobuf

This is just a simple personal reference for protobuf/gRPC.

### controller

Runs multiple time.tickers based on what is contained in the file under -job-list.  Based on the contents of the job list, the controller makes gRPC calls to two different services to drive them.  Currently the two services (cmd/svc[12]) are identified through command line flags (-svc-report & -svc-state).  It works, but would be better served with a service registry (e.g. svc[12] register with the registry on start/stop and the controller 'finds' them through the registry on each call - to-do #1).  It could also be inproved by sending the context cancel time to the server using timeduration.proto ( to-do #2 ).  This could all be extended for a large number of servers providing different services and to operate in K8s ( to-do #3 ).  The time.ticker  intervals are set for time.Second in the reference, but could be modified to hours or days.

#### job-list file

The job list is a tab delimited file with the following fields

    DEVICE = The name of the device
    JOB_TYPE = The type of job to execute (requires modification of the controller main.go RunJob() switch clause)
    INTERVAL = An integer that is multiplied by the time interval of the time.ticker

    EXAMPLE 

    DEVICE	JOB_TYPE	INTERVAL
    device1	state	60
    device2	state	60
    device3	state	120
    device4	state	30
    device5	state	60
    device5	report	24

### github.com/rpcox/grpc-protobuf/pkg/job

The controller and svc[12] use this package.

    go get -u github.com/rpcox/grpc-protobuf/pkg/job@latest

### svc[12]

Simple servers that receive a gRPC call and 'do something'.  Currently they just mimic and action through time.Sleep().