# DiningPhilosopher
Dining Philosopher - Mutual Exclusion by Message Passing

Machine 1: 140.158.131.163	: This is where I ran the Print Server that is connected through UDP connection by running the command 
"go run  Server.go 4444” 
where 4444 is the port number that I assigned manually while running. Once that’s done, it pops up the message that the server is ready in port 4444 which is going to be used while connecting all other philosophers/forks.

      Figure: 01_Server Initiation.PNG
Machine 2: 140.158.128.20: I assigned 3 forks and 2 philosophers in this machine. To do that, I opened 5 separate terminals. Used following commands to run Forks and Philosophers respectively.
* For F0: go run TCPFork.go 140.158.131.163:4444 140.158.128.20:4001
* For F1: go run TCPFork.go 140.158.131.163:4444 140.158.128.20:4002
* For F2: go run TCPFork.go 140.158.131.163:4444 140.158.128.20:4003
* For P0: go run TCPPhilosopher.go 140.158.131.163:4444
* For P1: go run TCPPhilosopher.go 140.158.131.163:4444

									                                           Figure: 02_Machine2 Initiation.PNG
Machine 3: 140.158.130.213	:  I assigned 2 forks and 3 philosophers in this machine. To do that, I opened 5 separate terminals. Used following commands to run Forks and Philosophers respectively.
* For F3: go run TCPFork.go 140.158.131.163:4444 140.158.130.213:4004
* For F4: go run TCPFork.go 140.158.131.163:4444 140.158.130.213:4005
* For P2: go run TCPPhilosopher.go 140.158.131.163:4444
* For P3: go run TCPPhilosopher.go 140.158.131.163:4444
* For P4: go run TCPPhilosopher.go 140.158.131.163:4444

Once all these setup is done, the process should start and the output should be seen like below:


