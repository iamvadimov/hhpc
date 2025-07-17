## Smart Energy Management System

Revision: 1.1

Date: June 30, 2025

Description: Building a monitoring subsystem running on low-powered 
hardware such as the Raspberry Pi to monitor the Smart Energy Management System.

### Prerequisites

While there are ready-to-use homes's energy efficiency solutions
available in the market, the goal of Smart Energy Management System is
to develop its own unique solution. Using Smart Energy Management System
applications, you can improve your home's energy efficiency. You can
optimise energy consumption tasks that are often forgotten and increase
your convenience and comfort.

Before you can start building the Smart Energy Management System
applications, you first need to collect the hardware and software
resources you'll use to construct the project. Then you have to create a
robust infrastructure to host this project, using the same scalable
tools that large companies and premier technology organisations use in
their own IT operational environments.

Our objective for the monitoring project as a part of Smart Energy
Management System applications was to avoid as much electrical
engineering and wiring as possible. The monitoring project was completed
without ever picking up a soldering gun. While it's commendable to use
one for appropriate cases, this part of project focuses more on software
than hardware. We also didn't want to have hardware components fail as a
result of poor soldering or confusing wiring diagrams, so we opted to
make the hardware configuration for this project as simple as possible
to avoid any frustration or expensive mistakes.

In this part of the Smart Energy Management System, we use a small
device, the Raspberry Pi, and the Go programming language to develop a
monitoring system.

This solution requires both hardware and software to make it fully
functional. Sensors and controllers need to connect to a computer
running code that knows how to interact with those devices. But at this
stage of the project development, there is no possibility to receive
real data from external energy consumers. In this situation, external
consumers are represented by a dataset.

### About Dataset

This dataset contains 525600 measurements gathered between January 2024
and December 2024 (12 months). A sample from the dataset is presented below

| Date      | Time    | Global_active_power | Global_reactive_power | Voltage | Global_intensity | Sub_metering_1 | Sub_metering_2 | Sub_metering_3 | 
| ----------| --------| -------------------:| ---------------------:| -------:| ----------------:| --------------:| --------------:| --------------:| 
| 14/12/2024| 17:58:00| 2.278               | 0.050                 | 236.740 | 9.600            | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 17:59:00| 2.466               | 0.092                 | 236.800 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:00:00| 2.464               | 0.092                 | 236.950 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:01:00| 2.470               | 0.094                 | 237.280 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:02:00| 2.482               | 0.096                 | 237.900 | 10.400           | 0.000          | 0.000          | 18.000         | 
| 14/12/2024| 18:03:00| 2.484               | 0.096                 | 237.820 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:04:00| 2.484               | 0.098                 | 237.960 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:05:00| 2.490               | 0.098                 | 238.480 | 10.400           | 0.000          | 0.000          | 18.000         | 
| 14/12/2024| 18:06:00| 2.492               | 0.098                 | 238.520 | 10.400           | 0.000          | 0.000          | 17.000         | 
| 14/12/2024| 18:07:00| 2.482               | 0.096                 | 238.120 | 10.400           | 0.000          | 0.000          | 18.000         | 


#### Data Set Information:

**Context:** Measurements of electric power consumption in one household
with a one-minute sampling rate over a period of 1 year. Different
electrical quantities and some sub-metering values are available.

**Data Set Characteristics:** Multivariate, Time-Series

**Attribute Information:**

1.  `date`: Date in format dd/mm/yyyy
2.  `time`: time in format hh:mm:ss
3.  `global_active_power`: household global minute-averaged active power
    (in kilowatt)
4.  `global_reactive_power`: household global minute-averaged reactive
    power (in kilowatt)
5.  `voltage`: minute-averaged voltage (in volt)
6.  `global_intensity`: household global minute-averaged current
    intensity (in ampere)
7.  `sub_metering_1`: energy sub-metering No. 1 (in watt-hour of active
    energy). It corresponds to the kitchen, containing mainly a
    dishwasher, an oven and a microwave (hot plates are not electric but
    gas powered).
8.  `sub_metering_2`: energy sub-metering No. 2 (in watt-hour of active
    energy). It corresponds to the laundry room, containing a
    washing-machine, a tumble-drier, a refrigerator and a light.
9.  `sub_metering_3`: energy sub-metering No. 3 (in watt-hour of active
    energy). It corresponds to an electric water-heater and an
    air-conditioner.


**Note**


1.  `(global_active_power*1000/60 - sub_metering_1 - sub_metering_2 - sub_metering_3)`
    represents the active energy consumed every minute (in watt hour) in
    the household by electrical equipment not measured in sub-meterings
    1, 2 and 3.
2.  To make it real the dataset contains some missing values in the
    measurements (nearly 1,25% of the rows). All calendar timestamps are
    present in the dataset but for some timestamps, the measurement
    values are missing: a missing value is represented by the absence of
    value between two consecutive semi-colon attribute separators.
3.  By visualizing the data, it'll be easy to spot trends and changes in
    active power, current intensity, active energy etc. as well as
    assign alerts when defined thresholds are exceeded. For example, you
    can configure Grafana to email you when your volume goes higher than
    250 vols.
4.  The detailed description of the original data can be 
    found [here.](http://archive.ics.uci.edu/ml/datasets/Individual+household+electric+power+consumption)


### Building a REST API Server

A primary technology standard that it well be call upon in this project
is Representational State Transfer, better known as REST. Rather than
re-inventing the wheel when it comes to passing state from one machine
to another, here is used a robust mechanism that works just as
effectively in small projects as it will in top-tier enterprise
applications.

REST packages information into JavaScript Notation, or JSON, format and
transfers those details using standard HTTP protocol. It's one of the
most prevalent ways to transfer meaningful information from one machine
to another, and it's popular with both consumer and enterprise web
applications. The formatted JSON data obtained from the hhpc-server will
then be consumed by the Prometheus instance on the Raspberry Pi server.
Hhpc is a REST server that will poll the onboard data and report those
values, formatted in a JSON payload that can be consumed for further
analysis. The values in this JSON will be converted into
Prometheus-friendly formatting using an exporter program to perform the
polling and conversion. But first, we need to get the HHPC data off the
onboard dataset, format it into a JSON-friendly format, and have an HTTP
server ready to accept new connections and deliver the JSON payload.
Hhpc data are measurements of electric power consumption in one
household with a one-minute sampling rate over a period of 1 year.
Different electrical quantities and some sub-metering values are
available. The measurements should be good enough for the monitoring
purposes.

#### Go Programming Language

The project works with any version of Go 1.20 or higher. At the time of
writing this project, Go 1.24 is generally available. Go supports
cross-compilation, allowing you to build an executable version of your
program using a different platform from the one running Go. So, you
don't necessarily need to run Go on your Raspberry Pi to complete this
project. If you want to install and use a different Go version than the
one shipped with your device, or if you're looking to install Go in a
different operating system, take a look at the [Go 
Downloads](https://go.dev/dl/) webpage.

    vadim@vadim:~$ go version
    go version go1.24.4 linux/amd64

#### Compiling REST API Server

We use Go to build a REST API server to pass messages and states between
the different components in your project. Running a Go application by
using the command `go run` is a nice and quick way to verify if the
application is working. It's particularly useful during development
time, but to run it in a production environment, you need to compile
your application into a binary file.

In Go, the compilation process is also known as "building", and you run
it by using the command `go build`. In its most basic form, running go
build builds your application into a single binary file compatible with
your current operating system and architecture. It also dynamically
links your application to the current system libraries. This information
is usually enough to run the application on your machine. However, using
this approach isn't the best way to run a Go application inside a
container.

    vadim@vadim:~$ cd hhpcserver/
    vadim@vadim:~/hhpcserver$ go build .
    vadim@vadim:~/hhpcserver$ ls -l
    -rw-r--r-- 1 vadim vadim       29 May 16 16:24 go.mod
    -rwxrwxr-x 1 vadim vadim  8532890 Jun 25 21:27 hhpcserver
    -rw-r--r-- 1 vadim vadim     1842 May 12 12:55 initdata.go
    -rw-r--r-- 1 vadim vadim     3251 May 12 21:20 main.go
    vadim@vadim:~/hhpcserver$ 

When building your Go application, you can specify additional
environment variables or parameters to customize the build and optimize
the resulting binary. Let's look at some common environment variables:

-   `GOOS`: Defines the target operating system for the build. Defaults
    to the current operating system. Go allows cross compilation, which
    means you can compile your Go application for a different operating
    system from your current system. To build the API server binary to
    run in a Linux container, set this variable to "linux".
-   `GOARCH`: Defines the target processor architecture for the build.
    Defaults to the CPU architecture of the system running the build. To
    build your API server binary to run on Raspberry Pi, set this
    variable to "arm64".
-   `CGO_ENABLED`: Defines whether or not the resulting binary
    dynamically links to the system libraries. To build a static binary
    to run in a container, set it to "0" (zero).

You can see all possible combinations of operating system and
architecture to use with `GOOS` and `GOARCH` by running go tool dist
list. In addition to variables, you can also pass additional parameters
to go build. To optimize your binary to run in a container and decrease
its size, pass the link parameter `-ldflags` set to `-s` `-w` to strip
the resulting binary of all debugging symbols. You can learn more about
additional build options by running go help build.

Build your API server binary with all these options by running go build
like this: :

    $ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w"

This command creates a binary file, restapi. Use the Linux command file
on it, to verify build details: :

    $ file restapi
    restapi: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, Go BuildID=QSMA....CAxpj, stripped

Now that the REST API server is ready, it can be run on Raspberry Pi.

### Running REST API Server as a Service

    vadim@vadim:~/hhpcserver$ ./hhpcserver
    Starting API server on port 4000...

    vadim@vadim:~$ date
    Wed Jun 25 23:32:57 EEST 2025

    vadim@vadim:~$ curl http://localhost:4000/api/v1/gethpc
    {"timestamp":"2025-06-25 23:33:00","gapower":0.288,"grpower":0.12,"voltage":242.41,"gintens":1.2,"sm1":2,"sm2":1,"sm3":0}

    vadim@vadim:~$ curl localhost:4000/api/v1/gethpc
    {"timestamp":"2025-06-25 23:35:00","gapower":0.284,"grpower":0.114,"voltage":242.66,"gintens":1.2,"sm1":1,"sm2":2,"sm3":0}

### Monitoring with Prometheus

Now that a a REST endpoint `http://localhost:4000/api/v1/gethpc` exists
to query at any time, use that always-on accessibility to your advantage
by leveraging the power of the [Prometheus](https://prometheus.io/)
server that will be installed a little bit later. Prometheus uses its
own particular naming conventions and formatting rules for indicating
value names and related assignments. As such, you'll need to convert the
JSON payload into an exported format that Prometheus can consume.

This project requires monitoring and alerting capabilities. Instead of
developing a custom solution to support this application only, let's
roll out an instance of Prometheus, a popular monitoring and alerting
solution, using containers on the Raspberry Pi. Prometheus is an open
source, metrics-based monitoring system.

Prometheus is a monitoring system and time-series database widely used
to monitor IT infrastructure, particularly container-based workloads.
Prometheus is developed using Go, and its flexibility and scalability
allow it to capture a variety of metrics, from small environments to
large corporate data centers. Prometheus works by "scraping" - or
polling on a fixed scheduled basis - systems for metrics and storing
them in a time-series database. It offers a powerful and fast query
language to retrieve and evaluate this data for correlation,
visualization, and alerting.

The target-monitored systems provide metrics to Prometheus via
"exporters." Prometheus ships several default exporters to monitor
common infrastructure such as operating systems and databases. It also
provides libraries for different programming languages to develop custom
exporters for your applications.

#### Compiling the Prometheus Exporter

This implementation, while not sophisticated, applies several of
[Prometheus exporter development
guidelines](https://www.prometheus.io/docs/instrumenting/writing_exporters/),
ensuring it exports metrics according to Prometheus's conventions, and
is safe for concurrent use. With all the code in place and the defined
REST server up and running, compile and run this Go Prometheus exporter
using the typical `go build` command:

    vadim@vadim:~ $ cd hhpcexporter/
    vadim@vadim:~/hhpcexporter$ go get github.com/prometheus/client_golang/prometheus
    go: downloading github.com/prometheus/client_golang v1.22.0
    go: downloading github.com/beorn7/perks v1.0.1
    go: downloading github.com/prometheus/common v0.62.0
    go: downloading github.com/cespare/xxhash/v2 v2.3.0
    go: downloading github.com/prometheus/client_model v0.6.1
    go: downloading github.com/prometheus/procfs v0.15.1
    go: downloading golang.org/x/sys v0.30.0
    go: downloading google.golang.org/protobuf v1.36.5
    go: downloading github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
    vadim@vadim:~/hhpcexporter$ go get github.com/prometheus/client_golang/prometheus
    vadim@vadim:~/hhpcexporter$ go get github.com/prometheus/client_golang/prometheus/promhttp
    vadim@vadim:~/hhpcexporter$ go build .
    vadim@vadim:~/hhpcexporter$ ls -l
    -rw-r--r-- 1 vadim vadim      512 May 16 15:31 go.mod
    -rw-r--r-- 1 vadim vadim     1629 May 16 15:31 go.sum
    -rwxrwxr-x 1 vadim vadim 13287297 Jun 26 21:30 hhpcexporter
    -rw-r--r-- 1 vadim vadim     5299 May 16 16:59 main.go
    -rw-r--r-- 1 vadim vadim      185 May 16 15:29 rootPage.html
    vadim@vadim:~/hhpcexporter$ 

#### Running the Prometheus Exporter

Set the environment variable `HHPC_SERVER_URL` to your REST API server
URL to configure the server to obtain `hhpc` data:

    vadim@vadim:~/hhpcexporter$ HHPC_SERVER_URL=http://localhost:4000/api/v1/gethpc ./hhpcexporter

Assuming no errors, you should be able to open a browser on your local
machine and visit `http://localhost:3030` with your preferred web
browser. If successful, you should see something similar to the response
in the next screenshot. Selecting the `/metrics` link on the page should
show the Prometheus-formatted results, including the defined `hhpc`
labels and associated values.

    vadim@vadim:~$ curl http://localhost:3030/metrics
    # HELP global_active_power Household global minute-averaged active power.
    # TYPE global_active_power gauge
    global_active_power{unit="kilowatt"} 0.602
    # HELP global_intensity Household global minute-averaged current intensity.
    # TYPE global_intensity gauge
    global_intensity{unit="ampere"} 2.8
    # HELP global_reactive_power Household global minute-averaged reactive power.
    # TYPE global_reactive_power gauge
    global_reactive_power{unit="kilowatt"} 0.092
    . . .
    # HELP hhpc_up Hhpc Server Status.
    # TYPE hhpc_up gauge
    hhpc_up 1
    . . .
    # HELP sub_metering_1 Energy sub-metering No. 1.
    # TYPE sub_metering_1 gauge
    sub_metering_1{unit="watt-hour of active energy"} 0
    # HELP sub_metering_2 Energy sub-metering No. 2.
    # TYPE sub_metering_2 gauge
    sub_metering_2{unit="watt-hour of active energy"} 0
    # HELP sub_metering_3 Energy sub-metering No. 3.
    # TYPE sub_metering_3 gauge
    sub_metering_3{unit="watt-hour of active energy"} 0
    # HELP voltage Minute-averaged voltage.
    # TYPE voltage gauge
    voltage{unit="volt"} 239.02
    vadim@vadim:~$

#### Install Prometheus Using Pre-compiled Binaries

Prometheus developers provide precompiled binaries for most official
Prometheus components. Check out the [download
section](https://prometheus.io/download/) for a list of all available
versions.

You can use the `curl` command to download
`prometheus-2.53.4.linux-amd64.tar.gz` file from a web server, then
extract it:

    vadim@vadim:~/$  % curl -L https://github.com/prometheus/prometheus/releases/download/v2.53.5/prometheus-2.53.5.linux-amd64.tar.gz | tar -xz

      % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                     Dload  Upload   Total   Spent    Left  Speed
      0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
    100  101M  100  101M    0     0  6978k      0  0:00:14  0:00:14 --:--:-- 7274k

    vadim@vadim:~/$ ls -l pro*
    drwxr-xr-x 4 vadim vadim      4096 Mar 18 17:08 prometheus-2.53.4.linux-amd64

The Prometheus server is a single binary called `prometheus`. You can
run the binary and see help on its options by passing the `--help` flag.

    vadim@vadim:~/$ cd prometheus-2.53.4.linux-amd64/
    vadim@vadim:~/prometheus-2.53.4.linux-amd64$ ./prometheus --help
    usage: prometheus [<flags>]

    The Prometheus monitoring server

    . . .

#### Configuring Prometheus to Query the Exporter

However, we also want to deploy a second application running a default
`node-exporter` to collect and expose data from the Linux OS on your
Raspberry Pi host. You can use the `curl` command to download
`node_exporter` file from a web server, then extract it:

    vadim@vadim:~/$ % curl -L https://github.com/prometheus/node_exporter/releases/download/v1.9.1/node_exporter-1.9.1.linux-amd64.tar.gz | tar -xz
      % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                     Dload  Upload   Total   Spent    Left  Speed
      0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
    100 11.0M  100 11.0M    0     0  4618k      0  0:00:02  0:00:02 --:--:-- 6892k
    vadim@vadim:~/$ cd node_exporter-1.9.1.linux-amd64
    vadim@vadim:~/node_exporter-1.9.1.linux-amd64$

Before starting Prometheus, let\'s configure it. First, create a custom
Prometheus configuration file. Most of these configurations are
standard, and you can find a sample file and detailed explanation about
this configuration in the project's [documentation
page](https://prometheus.io/docs/prometheus/latest/configuration/configuration/).
By default, your Prometheus instance will scrape the target system for
metrics every fifteen seconds. Change to a different value if you want
it to be more or less aggressive. In the `scrape_configs` section, you
have the default configuration to collect metrics from itself as
`job_name: prometheus`.

Prometheus configuration is [YAML](https://yaml.org/). The Prometheus
download comes with a sample configuration in a file called
`prometheus.yml` that is a good place to get started.

    # my global config
    global:
      scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
      evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
      # scrape_timeout is set to the global default (10s).

    # Alertmanager configuration
    alerting:
      alertmanagers:
        - static_configs:
            - targets:
              # - alertmanager:9093

    # Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
    rule_files:
      # - "first_rules.yml"
      # - "second_rules.yml"

    # A scrape configuration containing exactly one endpoint to scrape:
    # Here it's Prometheus itself.
    scrape_configs:
      # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
      - job_name: "prometheus"

        # metrics_path defaults to '/metrics'
        # scheme defaults to 'http'.
        static_configs:
          - targets: ["localhost:9090"]

      - job_name: "node"
        static_configs:
          - targets: ["localhost:9100"]

      - job_name: "hhpc"
        static_configs:
          - targets: ["localhost:3030"]

#### Starting Project

    vadim@vadim:~/$ cd hhpcserver/
    vadim@vadim:~/hhpcserver$ ./hhpcserver
    Starting API server on port 4000...
    Processing request from 127.0.0.1:49540: User-Agent: Go-http-client/1.1

    vadim@vadim:~/$ cd hhpcexporter/
    vadim@vadim:~/hhpcexporter$ HHPC_SERVER_URL=http://localhost:4000/api/v1/gethpc ./hhpcexporter

    vadim@vadim:~/$ cd node_exporter-1.9.1.linux-amd64/
    vadim@vadim:~/node_exporter-1.9.1.linux-amd64$ ./node_exporter

    vadim@vadim:~/$ cd prometheus-2.53.4.linux-amd64/
    vadim@vadim:~/prometheus-2.53.4.linux-amd64$ ./prometheus

The project is suited to run in your browser instantly without any extra
configuration or downloading a software client:

    http://localhost:9090/

### Summary

In this guide, you installed `hhpc` applications and Prometheus,
configured a Prometheus instance to monitor resources. To continue
learning about Prometheus, check out the
[Overview](https://prometheus.io/docs/introduction/overview/) for some
ideas about what to explore next.

Now you know what the needs of the project are and what data needs to be
captured and evaluated. If you want to capture and manipulate more
details, that's your choice, not the choice of the equipment
manufacturer.
