# DNS (Drone Navigation Service)

## How to use

- `make build` - builds docker image
- `make run port=<int> sector_id=<int>` - runs docker container with exposed port and defined sector_id
- `make stop sector_id=<int>` - stop and remove container with such a sector_id
- `make thanos_snap` - removes all containers and images for DNS
- `make test` - run tests

`http://localhost:<port>/api/v1/location` - endpoint
`http://localhost:<port>/swagger/index.html` - openapi documentation

## Decisions made

- in task description request is json, but I decided to make pass data as query parameters,
because it's idiomatic GET request
- using decimal instead of float. I think that accuracy is important for coordinates handling
and losing accuracy because of float representation is unacceptable.
- sectorID is passed to docker container as env variable: with such a configuration we can easily 
run same docker image for different sectors

## Limitations

- Location calculation is a simple stub, it calculates meaningless number, but that's the task
- I think that something like JWT-auth is suitable here, without authentication anyone can get
location of data storages. I'm not sure, that this information should be public.
- No drone management. Drones should be able to register to DNS, receive jwt token for further communication.
Also, it's required in order to determine if drone from invalid sector (drone is in wrong sector because of mistake)
and reject exposing datacenter location (this information is for this sector drones only)
- Also, drone can be directed to a right sector, if DNS can communicate with each other (no hard code, 
with service mesh for example). Instead of rejecting drone request - it's redirected to the correct DNS for this drone,
correct coordinates received and drone moves to correct position.

## About layout

- api - swagger specification
- build - all for docker
- cmd - entrypoints to application, binaries are compiled from here
- internal - application code, forbidden fir importing by other packages
    - location - location domain, if we'll need drone management we'll create next to it a drone
    directory. This directory will contain all drone relaret code: handlers, dto, entities, services etc.
        - dto
        - handler
        - service
    - util - maybe this should be named as pkg, this package can be imported by all other internal packages
    
## Questions

**1) What instrumentation this service would need to ensure its observability and operational transparency?**

Hmm, about operational transparency- that's an interesting question. I'm not sure for this case, but maybe openapi
documentation with an interactive playground (swagger) makes the process if communication between the client and service more
transparent?
Observability is a must-have nowadays and in the future (I think k8s is dead there and smth much more powerful 
and much more complex exists) it's a must-have thing too. Observability mainly is about metrics (especially business), logs,
distributed traces. Amazon is definitely alive in the future and we can use Amazon CloudWatch if we use AWS.
If not, we can use any other solution, there are lots of them on the market. That's a question more about tools, I guess.
I used Thundra.io for AWS Lambda observability for example, it was much more powerful than AWS X-Ray and Cloudwatch.

**2) Why throttling is useful (if it is)? How would you implement it here?**

Throttling definitely has its usecases. For example, ddos protection: it's better to cut off requests and protect service
form being working more and more slowly and finally going out of, for example, memory. Also, there's another pretty common 
usecase: protect from resource over-usage. For example in many paid services user can define rate limiting. If user's
service is under ddos and user service sends tones of request to sms provider, only N req/min (for example) will be 
executed and this will save money of the user.
This can be implemented on api and application level. On app level in this task throttling can be implemented without doing
anything: chi router has built in throttle middleware. On api level, API Gateway can be used for this. Also in case of
running out of limit, API Gateway can serve cashed responses if this is acceptable (for GET request).


**3) What we have to change to make DNS be able to service several sectors at the same
   time?**
   
First of all, we need pass list of domain ids. This can be done in the same way with env variables, for example:
SECTOR_IDS=111,222,333. This string is `split(",")` at server side. Then, we need somehow define which sectorId
use for incoming request. If drone registration is implemented - we can store sector of the drone in his JWT token.
It's safe and can't be rewritten by Somalis Space Pirates wanting to rob the data storage (they will have no possibility
to find out the location).
   
**4) Our CEO wants to establish B2B integration with Mom's Friendly Robot Company by
   allowing cargo ships of MomCorp to use DNS. The only issue is - MomCorp software expects loc value 
   in location field, but math stays the same. How would you approach this? What’s would be your implementation
   strategy?**
   
It depends on the complexity of the solution and some factors. I see these possibilities:
- If drone registration is enabled, we can store in JWT the owner of the drone. For GetLocationRequest we can write
presenters: DefaultPresenter - for our drones, MomCormPresenter - for MomCorp. Pros: easy to implement. Cons: not very
scalable if we'll have lots of integrations.
- API Gateway. With this pattern implementation we can change requests and responses. Again, we must somehow define
the owner of the drone (jwt, ip, query param etc), but all modifications are handled outside of the app.
Pros: scalable, code is cleaner. Cons: more complex architecture of the app, more support dur to this increased complexity.
Slower to implement.
   
**5) Atlas Corp mathematicians made another breakthrough and now our navigation math is even better and more accurate, so 
   we started producing a new drone model, based on new math. How would you enable scenario where DNS can serve 
   both types of clients?**
   
This is taken into account in my implementation. LocationCalculator is an abstraction, and any implementation can be used.
LocationCalculator abstraction is injected to LocationService at bootstrap stage. So I see these options:
- different sectors can deploy DNS with different calculation strategy. It's suitable if drone models are launched sector by sector.
- if in one sector new and old drones are presented, we can create new api version, v2. Old drones will send requests
to /api/v1 and the new ones to /api/v2. Handlers for these endpoints will use different interface implementations.
- if drone can send his software version while registration to DNS we can store it in his JWT and just create a factory
for getting LocationCalculator implementaion. In this case implementation will be selected at runtime, and not at server launch.

**6) In general, how would you separate technical decision to deploy something from business decision 
   to release something?**
   
Deployment is technical decision. Release - business decision. I'm not a fan of "deploy and release at the same time" approach.
This approach leads to different problems and release rollback as a result. Deploying before release allows us to
see how new features work in production environment without high load to these features (also, users can event not be informed, these
new features are hidden for them). We can use A/B testing, canary release. For example, we can allow access to new feature for
1% of users (or even 1% of subset of users, only mac-users for example). If everything goes fine, we can increase this amount (and repeat until 100%). If something goes wring, we can just switch
of new feature, redirect all users to old service instances and fix all the problems.
