import argon2 from "argon2";
import type { Db } from "mongodb";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";

export const post: APIRoute = async ({ request }) => {
  try {
    const { username, email, firstName, lastName, password } =
      await request.json();

    const db: Db = await getDb();
    const users = db.collection("users");
    const registeredUser = await users.findOne({
      $or: [{ username: { $eq: username } }, { email: { $eq: email } }],
    });
    if (registeredUser) {
      return new Response(
        JSON.stringify({
          status: "error",
          message:
            "User already registered, please use other email or username!",
        }),
        {
          status: 422,
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }

    const passwordHash = await argon2.hash(password);
    const user = {
      firstName,
      lastName,
      username,
      email,
      password: passwordHash,
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    const insertedUser = await users.insertOne(user);

    return new Response(
      JSON.stringify({
        status: "success",
        data: {
          id: insertedUser.insertedId,
          firstName,
          lastName,
          username,
          email,
          createdAt: user.createdAt,
          updatedAt: user.updatedAt,
        },
      }),
      {
        status: 200,
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
  } catch (err) {
    console.error(err);
    return new Response(
      JSON.stringify({
        status: "error",
        message: "Your request can't be processed, please try again",
      }),
      {
        status: 422,
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
  }
};
