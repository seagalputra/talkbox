const { getObjectId } = require("mongo-seeding");

const messages = [
  {
    id: getObjectId("message1"),
    body: "Hello, how are you?",
    attachment: "",
    userId: getObjectId("john_doe"),
    roomId: getObjectId("room1"),
    createdAt: new Date(),
    updatedAt: new Date(),
  },
];

module.exports = messages;
