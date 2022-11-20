import bcrypt from "bcryptjs";
import Joi from "joi";
import type { Db } from "mongodb";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";
import { errorResponse, successResponse } from "../../../utils/api";

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

    if (password !== passwordConfirmation)
      return errorResponse(
        null,
        422,
        "Your password doesn't match, please re-check your password again"
      );

    const db: Db = await getDb();
    const users = db.collection("users");
    const registeredUser = await users.findOne({
      $or: [{ username: { $eq: username } }, { email: { $eq: email } }],
    });
    if (registeredUser)
      return errorResponse(
        null,
        422,
        "User already registered, please use other email or username!"
      );

    const passwordHash = await bcrypt.hash(password, 10)
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

    return successResponse({
      id: insertedUser.insertedId,
      firstName,
      lastName,
      username,
      email,
      createdAt: user.createdAt,
      updatedAt: user.updatedAt,
    });
  } catch (err) {
    console.error(err);
    return errorResponse(err);
  }
};
