#!/usr/bin/env node

db.createUser({
  user: "mongo",
  pwd: "mongo",
  roles: [
    {
      role: "readWrite",
      db: "talkbox",
    },
  ],
});
