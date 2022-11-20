import bcrypt from 'bcryptjs'
import Joi from "joi";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";
import { errorResponse, successResponse } from "../../../utils/api";
import { generateAuthToken } from "../../../utils/token";

// username can be filled with email too
const schema = Joi.object({
  username: Joi.string().required(),
  password: Joi.string().required(),
});

export const post: APIRoute = async ({ params, request }) => {
  try {
    const body = await request.json();

    const { username, password } = await schema.validateAsync(body, {
      abortEarly: false,
      errors: { escapeHtml: true },
    });

    const db = await getDb();
    const users = db.collection("users");
    const user = await users.findOne({
      $or: [
        {
          username,
        },
        {
          email: username,
        },
      ],
    });

    if (!user)
      return errorResponse(
        null,
        404,
        "Your account doesn't exists, please check your username again"
      );

    const validPassword = await bcrypt.compare(password, user?.password);
    if (!validPassword)
      return errorResponse(
        null,
        422,
        "Your password is invalid, please check your password again"
      );

    const token = await generateAuthToken(user);

    return successResponse({
      id: user._id,
      firstName: user.firstName,
      lastName: user.lastName,
      email: user.email,
      username: user.username,
      authToken: token,
    });
  } catch (err) {
    console.error(err);
    return errorResponse(err);
  }
};
