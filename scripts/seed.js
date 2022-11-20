import { MongoClient } from "mongodb";
import argon2 from "argon2";
import dotenv from "dotenv";
dotenv.config();

const connectionUrl = process.env.DATABASE_URL;

const client = new MongoClient(connectionUrl);

async function run() {
  try {
    const db = client.db(process.env.DATABASE_NAME);
    const users = db.collection("users");
    const passwordHash = await argon2.hash("password");
    const john = {
      firstName: "John",
      lastName: "Doe",
      username: "johndoe",
      email: "john@email.com",
      avatar: "",
      password: passwordHash,
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    const johnQuery = { email: "john@email.com" };
    const johnUpdate = {
      $set: john,
    };
    await users.updateOne(johnQuery, johnUpdate, {
      upsert: true,
    });

    const alice = {
      firstName: "Alice",
      lastName: "Wonderland",
      username: "alice",
      email: "alice@email.com",
      avatar: "",
      password: passwordHash,
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    const aliceQuery = { email: "alice@email.com" };
    const aliceUpdate = {
      $set: alice,
    };
    await users.updateOne(aliceQuery, aliceUpdate, {
      upsert: true,
    });

    const rooms = db.collection("rooms");
    const roomQuery = {
      $and: [
        { "participants.username": john.username },
        { "participants.username": alice.username },
      ],
    };
    const roomUpdate = {
      $set: {
        participants: [
          {
            username: john.username,
            email: john.email,
          },
          {
            username: alice.username,
            email: alice.email,
          },
        ],
        roomType: "private",
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    };
    await rooms.updateOne(roomQuery, roomUpdate, { upsert: true });

    const johnResult = await users.findOne({ email: "john@email.com" });
    const roomResult = await rooms.findOne(roomQuery);
    const messages = db.collection("messages");
    const messageQuery = {
      userId: johnResult._id,
    };
    const messageUpdate = {
      $set: {
        body: "Hello, World!",
        attachment: "",
        userId: johnResult._id,
        roomId: roomResult._id,
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    };
    await messages.updateOne(messageQuery, messageUpdate, { upsert: true });
  } finally {
    await client.close();
  }
}

run().catch((err) => console.err(err));
