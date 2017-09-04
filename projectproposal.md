# Up Project Proposal

## Daniel Wilkins

### Definitions
There are a couple terms defined by the community that I'll be making use of here, some of which have their names for historical reasons; this is mainly so that if I 'accidentally' use the community terms my meaning is clear:
Dent: A message sent via OStatus. Derived from the original (now defunct) OStatus instance, iDENTi.ca.
Federation: Communication with the other servers in the 'fediverse'(below) so that it (mostly) does not matter whether everyone is on a single instance or multiple. Essentially the exchanging of posts and other status updates between servers.
Fediverse: All of the OStatus servers on the web, regardless of their specific implementation, be it GNU Social, Mastodon, Postactiv, etc.
MUST, SHOULD, MAY: borrowed from RFCs; MUST is required, SHOULD is recommended, MAY is optional. Think of these as categories A, B, and C.
Social: [GNU Social](http://gnu.io)


### Project Overview
Up will be a simple, lightweight, single-user server which implements the OStatus protocol for federated social networking. Up targets small instances where a person wants to be able to host it on whatever spare hardware they have on hand. Essentially, people like myself. As far as the scope of Capstone, I am aiming for essentially the least common denominator of OStatus support, what can be expected to federate easily no matter the software. 

#### Up MUST implement the following:
* The ability to authenticate. Given the single-user nature and the time constraints I may simply make this a configurable secret key that the server checks for on commands which require authentication. Because these instances are by definition single-user, something as simple as an oauth-style scheme where the server generates a 'key' that acts as a shared secret and the user can ask for a new key would be sufficient. Test: try to send a dent without authenticating, it should fail. 

* A REST API either of my own design or compatible with another OStatus implementation such as Social or Mastodon, depending on the complexity and how well the APIs match with Up's goals. Test: able to interact with the site(?)

* The ability to submit text dents, which are then stored and distributed as specified below. Test: post a dent, see if it appears on the timeline of a test account on another instance.

* Permalinks for every text dent submitted. These could be as simple as an incrementing "dent ID" as long as it is easily accessible.

* The ability to subscribe and subscribe to, that is, to support the subscription status as described in the current OStatus RFC. This means that if someone is subscribed to me, Up will send my posts to their server to be served to them, and if I subscribe to someone, Up will accept that person's posts from their server and serve them to me. This is the tricky part. Test: subscribe to a user on a test instance, have that user subscribe to mine, test communication (does a text post get back and forth?)

* Have enough of a web interface, no matter how ugly, to accomplish the above and to provide permalinks for dents. Test: Implied above: successfully be able to subscribe.

* The ability to direct a dent to someone, even if they are not subscribed, with @person@instance. Test: Set up another test account with no subscriber/subscribee relationship, send messages, see if they get through.

### Up SHALL NOT have the following: 
* Configuration via web interface. Because the target of this is a single-user, self-hosted instance, config files in some well understood format (json, ini, etc.) should be sufficient.
* Multiple users. Small-scale multiple user systems would be a goal for a successor project, not for this one. Federation is hard enough with the assumption that our instance is already single-user.

#### Up SHOULD implement the following:
* Be easy to deploy, by my definition of easy to deploy: it should behave well on a server which is used by other programs (Pump.io, another program in the same space but using an incompatible protocol, famously did not like sharing a web server with other services,) and not use excessive resources. Additionally, mucking with a database should not be needed: it should use sqlite or at the very least automate as much database setup as possible. My personal instance of GNU Social is claiming about 512MB between the programs and the database; Mastodon wants much more. I would like to beat GNU Social on memory usage. Test: set up an instance, subscribe roughly on the order that I do on other instances, look at memory usage.

* Have a client. While in the worst case Up could be driven by, say, cURL or netcat, it would be better to either be compatible with an existing client if API compatibility is practical, or have a simple client if not. Test: be able to use the software without just posting through cURL.

* Support 'hashtags'. A message can be 'tagged' by adding a string of the form '#tag' to it. This is then stored in the database under that tag and can be found by asking the server for posts tagged with 'tag'. Test: Hashtags are *not* distributed, and so posting based on that and being able to search for tags posted by it through the API should be sufficient.

* HTTPS. As something of a proof of concept, HTTPS could conceivably be saved until later, but in order for it to be usable and replace my current setup by the end of the semester, it would be needed. The only problem I can see which would prevent this is web server configuration problems.

#### Up MAY implement the following:
* A full web interface. I generally prefer to use clients, but I imagine that some people prefer interfaces they can access from anywhere.

* Blocking. The ability to refuse to accept messages from a certain person, so that even if a dent is directed towards me I will not receive it. Test: Block a test account, see if direct messages get through (or groups, if groups are implemented)

* 'Pretty' directing. The ability to simply do @person rather than @person@instance if the name is unambiguous (I am only subscribing to one @person). Test: Subscribe to an unambiguous user off-instance, try to direct message.

* Groups. Mastodon does not implement this, but on Social and other implementations such as Postactiv one can subscribe to a group, which is a federated set of users with a common interest. I can send a dent to a group by including '!group' in the name and this message is sent to all subscribers of that group regardless of whether they are subscribed to me. Test: Subscribe to a group on the fediverse (!coffee is a good one), see if posts percolate through.

### Similar Work
The protocol has been implemented by various programs, dedicated to the protocol and not. Most notably [GNU Social](http://gnu.io), [Mediagoblin](http://mediagoblin.org), and [Mastodon](http://mastodon.social). It has also been implemented in programs such as [Friendica](http://friendi.ca). The primary differences between Up and these programs is simplicity: it should be both easy to deploy (ideally you build the program, do something like ```./up init && ./up start``` and then you tell your web server to proxy it).

### Previous Experience
I have worked on several web apps for my own use: pastebins, image galleries, wikis, and the like. Similarly I've dealt with XML and JSON parsing in, for example, my podcast fetcher. Finally, I've administered GNU Social instances since before it was called GNU Social, probably the past 7 years or so; it was one of the first things I got running when I got a server.

### Technology
* Go: the only real risk factor in terms of technology, but given that I've been following Go since its announcement and have played around with its predecessor, Limbo, I am fairly confident that I can pick it up quickly.

* SQLite: For the database and anything which won't fit into a config file. Needed to store messages, follows, following, etc. Handled by Mattn's [go-sqlite3](http://github.com/mattn/go-sqlite3) library. By all means it looks like SQLite bindings for any other language, just in Go.

* OStatus: This implies various other (but fairly simple) protocols such as Salmon and PuSH (although the standard does make Salmon optional). The individual standards themselves appear fairly simple, however, with some essentially just defining an xml schema or some other file format (like webfinger). Essentially it's a lot of XML over HTTP, both of which Go supports in the standard library. The RFC specifically states that the authors hope that some parts of the specification are such obvious applications that anyone familiar with the specifications in question would put them together in such a way.

* XML: I'm not exactly sure that this qualifies as a technology, but most of the standards in OStatus involve XML in one way or another; I will be reading an awful lot of dumps in that format. XML support is in the Go standard library.

* Webfinger: A simple protocol that defines a special URL and a JSON schema, essentially; webfinger.net recommends [ant0ine's go-webfinger](http://github.com/ant0ine/go-webfinger) for client work; in this context on the server side the webfinger information would be a static file to be served.

### Risk

* (Low) Go doesn't work for me: I doubt this will be the case but I should be able to move to another language I'm more familiar with without trouble, even if it won't be quite so efficient. My specific backup is Perl and Mojolicious, a web framework I am at least somewhat familiar with.

 * (High) Standard Drift: Many of the RFCs listed in the OStatus specification were hosted on Google Code, but are accessible from Wayback Machine. It may be that implementations have drifted considerably from the standard, and so OStatus would now be a de facto standard with a lot of internal knowledge needed. Thankfully it is fairly easy to contact those who are implementing the protocol these days via OStatus itself and via IRC, not to mention email. I know of at least one other person who is creating their own OStatus server right now, for example. To caution against this, I would like to get federation working as quickly as possible, at least in a read-only form; this way I have the maximum amount of time to squish bugs in it. If I can't find what the de facto standard of the day is despite my attempts to mitigate the danger, then the project is simply dead in the water.
 
 * (Low) Lightning strikes my server. Up would be running on the server in my guest bedroom and if there were a catastrophic hardware failure I would need to try to move to AWS or similar, which would take time becuase I will mostly be working with real world data. If there *is* a catastrophic failure, I would probably move to AWS or another hosting service which gives students free access and create an account on a public instance to create test data.

* (Medium) API troubles. While I have worked on web development projects before, I've never actually had to design a REST API before. If I decide that existing APIs are a bad match for Up there will be some amount of friction as I learn.

* (Medium) Proxying Problems. I haven't done proxying to an application running on HTTPS in quite a while, and I think that a proxied solution rather than FastCGI is the default for Go. There will be some configuration tweaking needed to get it working properly, and if I can't get it working properly then I'll have to figure something else out (Run only Up?)
