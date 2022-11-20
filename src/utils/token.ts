import jwt, { JwtPayload } from "jsonwebtoken";
import { config } from "../assets/app";

function generateAuthToken(user: any): Promise<string | undefined> {
  return new Promise((resolve, reject) => {
    const secret: string = config.jwtSecret ?? "secret";
    jwt.sign(
      {
        id: user._id,
        firstName: user.firstName,
        lastName: user.lastName,
        username: user.username,
        email: user.email,
      },
      secret,
      { expiresIn: "90d" },
      function (err, token) {
        if (err) {
          reject(err);
        }

        resolve(token);
      }
    );
  });
}

function validateAuthToken(
  token: string
): Promise<string | JwtPayload | undefined> {
  return new Promise((resolve, reject) => {
    const secret: string = config.jwtSecret ?? "secret";
    jwt.verify(token, secret, function (err, decoded) {
      if (err) {
        reject(err);
      }

      resolve(decoded);
    });
  });
}

export { generateAuthToken, validateAuthToken };
