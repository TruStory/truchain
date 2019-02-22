Notes Push Notif Architecture

Option 1: Centralized

Client calls truapi POST /device 
- maintain table to addresses + device tokens

Push notif Table
- maintain queue of undelivered push_notifs in postgres DB

Push Notif Deliverer
- separate go function  popping notifs off queue and delivering them. Started and stopped alongside  web server.


Option 2: Dencentralized "Notifier Micro Service"

Client calls API of serparate web server POST /device token

Micro service maintains table of addresses + device tokens

Micro service subscribes to truchaind websocket API

Micro service contains business logic of what notifs to send to who upon which events

UNKNOWNS: Are Tendermint/Cosmos Websocket events  easily augmented with Msg specific fields? Like Challenge.story.creator
WORKAROUND: Micro service could call truchaind /graphiql for additional required data

BENEFTS: 
  No duplicate notifications
  Allows scaling/decentralization
  Reduces complexity of root node
  


Design Pattern Common to both Options:

QUEUE!
Schedule notif by pushing onto queue
Deliver notif at top of queue
Update notif as delivered.

Allows for "dashboard" style insights into backlog of notifs, errors, etc.

Allows for crash-recovery, should deliverer process go down. 

Allows for "mass notifs" such as notifying every user upon new Story creation.
Incremental queue-based approach delivers all notifs eventually  without crashing the full node or notif deliverer process.







