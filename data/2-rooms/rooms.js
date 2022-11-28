const { getObjectId } = require("mongo-seeding");

const rooms = [
  {
    id: getObjectId("room1"),
    participants: [
      {
        id: getObjectId("john_doe"),
        username: "john_doe",
        email: "johndoe@email.com",
      },
      {
        id: getObjectId("alex_roger"),
        username: "alex_roger",
        email: "alexroger@email.com",
      },
    ],
    roomType: "private",
    createdAt: new Date(),
    updatedAt: new Date(),
  },
];

module.exports = rooms;
