const { getObjectId } = require("mongo-seeding");
const bcrypt = require("bcryptjs");

const passwordHash = bcrypt.hashSync("password", 10);

const users = [
  {
    firstName: "John",
    lastName: "Doe",
    username: "john_doe",
    email: "johndoe@email.com",
    avatar:
      "https://i.picsum.photos/id/524/200/200.jpg?hmac=t6LNfKKZ41wUVh8ktcFHag3CGQDzovGpZquMO5cbH-o",
    password: passwordHash,
  },
  {
    firstName: "Alex",
    lastName: "Roger",
    username: "alex_roger",
    email: "alexroger@email.com",
    avatar:
      "https://i.picsum.photos/id/213/200/200.jpg?hmac=Jzh2fbzIE1nc6J8qLi_ljVCRz0AITXxCC1Z8t2sD4jU",
    password: passwordHash,
  },
];

module.exports = users.map((user) => {
  return {
    ...user,
    id: getObjectId(user?.username),
    status: "active",
    createdAt: new Date(),
    updatedAt: new Date(),
  };
});
