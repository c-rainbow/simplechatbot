package spamcheck

/*
Plugin to check for spam messages in chat
It listens to all messages, and if a chat message looks like a chat, it either
(1) Warns the chatter
(2) Timeout the chatter for N seconds (N can be different, depends on previous timeout history)
    i) Also warns the chatter
(3) Permanently ban the chatter
    i) Also say something cool to the banned chatter
*/
