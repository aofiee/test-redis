Installation
1. create "db" directory in root project for store postgres database
2. config .env in service1 and service2 project

service1 env
```
APP_PORT=9090
POSTGRES_HOST=database
POSTGRES_PORT=5432
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=postgres
POSTGRES_SSLMODE=false

REDIS_PASSWORD = password
REDIS_HOST = redis
REDIS_PORT = 6379
```

service2 env
```
APP_PORT=9091
POSTGRES_HOST=database
POSTGRES_PORT=5432
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=postgres
POSTGRES_SSLMODE=false

REDIS_PASSWORD = password
REDIS_HOST = redis
REDIS_PORT = 6379
```

3. run docker to start database, redis, go service
```
docker-compose up -d && docker-compose logs --follow
```
4. Have a document for rest api in service1 and service2 directory "docs.rest"
use in vscode extension "rest client" for create machine and stock

*** Note
In service2 on branch "unit-test" I code some unit test
can run in directory service2 with command
```
go test ./... -cover 
```

Best Regards,.
Khomkrid L.

Instructions
1. Build the Publish-Subscribe mechanism. Allow ISubscriber objects to register against an concrete IPublishSubscribeService object for an event type. Implement the publish method so that when a publish event occurs, all subscribers of that the event type published will have a chance to handle the event. The subscribers should be working off a shared array of Machine objects, mutating them depending on the event received.
2. Now add the method 'unsubscribe' on IPublishSubscribeService to allow handlers to unsubscribe from events. You may change the existing method signatures.
3. Implement MachineRefillSubscriber. It will increase the stock quantity of the machine.
4. If a machine stock levels drops below 3 a new Event, LowStockWarningEvent should fire. When the stock levels hits 3 or above (because of a MachineRefillEvent), a StockLevelOkEvent should fire. You may want to introduce new subscribers (e.g. a new subscriber called StockWarningSubscriber). For each machine,  LowStockWarningEvent or StockLevelOkEvent should only fire one time when crossing the threshold of 3. Remember subscribers should be notified in the order of the events occured.

Your program should now allow you to create ISubscriber objects, register them using your IPublishSubscribService implementation. You can then create IEvent objects and call your IPublishSubscribService's implementations .publish() method. All handlers subscribed should have their 'handle' methods invoked.

Note I: Handlers can also create new events, if desired. The events would get handled after all existing events are handled.

Note II: If a subscriber subscribes after an event has already been published and consumed, they will not receive that event.

You may make any changes to this codebase as long as you ultimately build a Pub-Sub application capable of handling the existing machine sale and refill events.

Please share your work using a GitHub repository with @proftom.
