import { ObjectId } from "mongodb";
import { validateAuthToken } from "./token";
import { getDb } from "../assets/mongo";

async function authenticateUser(
  request: Request
): Promise<{ user?: any; isAuthorized: boolean }> {
  try {
    const authToken = request.headers.get("Authorization");
    const [, plainToken] = authToken?.split(" ") || ["Bearer"];
    if (!plainToken) {
      return { isAuthorized: false };
    }

    const user = await validateAuthToken(plainToken);
    if (!user) {
      return {
        isAuthorized: false,
      };
    }

    const db = await getDb();
    const users = db.collection("users");
    const availableUser = await users.findOne({ _id: new ObjectId(user.id) });
    if (!availableUser) {
      return {
        isAuthorized: false,
      };
    }

    return {
      user,
      isAuthorized: true,
    };
  } catch (err) {
    console.error(err);
    return {
      isAuthorized: false,
    };
  }
}

export default authenticateUser;
