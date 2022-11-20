import argon2 from "argon2";
import Joi from "joi";
import type { Db } from "mongodb";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";

const schema = Joi.object({
  firstName: Joi.string().optional(),
  lastName: Joi.string().optional(),
  username: Joi.string().alphanum().min(3).max(30).required(),
  password: Joi.string().pattern(new RegExp("^[a-zA-Z0-9]{3,30}$")).required(),
  passwordConfirmation: Joi.ref("password"),
  email: Joi.string().email({ minDomainSegments: 2 }).required(),
});

export const post: APIRoute = async ({ request }) => {
  try {
    const body = await request.json();

    const {
      firstName,
      lastName,
      username,
      password,
      passwordConfirmation,
      email,
    } = await schema.validateAsync(body, {
      abortEarly: false,
      errors: { escapeHtml: true },
    });

    if (password !== passwordConfirmation) {
      return new Response(
        JSON.stringify({
          status: "error",
          message:
            "Your password doesn't match, please re-check your password again",
        }),
        {
          status: 422,
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }

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
    if (err instanceof Joi.ValidationError) {
      const errorDetails = err.details.map((detail) => ({
        field: detail?.path[0],
        message: detail?.message,
      }));

      return new Response(
        JSON.stringify({
          status: "error",
          message: errorDetails,
        }),
        {
          status: 400,
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }

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
