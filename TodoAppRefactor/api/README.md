## Handlers

Contains the handlers for each api request (Create, Get, GetAll, MarkAsComplete, Delete). Each handler has its own channel
which it can submit commands to. These commands are then processed through a 'RequestHandler' which directs these commands to 
the data service. 

Within the data service, read operations use 'RLock' whilst write operations use 'Lock'. The intention with
this is to have it so multiple read requests can happen at once, speeding up processing of requests.

## Contracts

The purpose of the contracts is to make sure that the data being passed into any requests are consistent. For example, when
making a new todo item, the 'CreateContract' is used to pass in the correct information.

## Mocks

Mock data service used for unit testing the API handlers.

## Responses

Responses are used to make sure the responses returned are consistent.

## Benchmarks

This is a benchmark test that sets up a server and performs a certain amount of random api calls with random parameters. A file
is also included to show example results from the test.
