import type { Db } from "mongodb";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";
import { getQueryParams } from "../../../utils/queryParams";
import { errorResponse, successResponse } from "../../../utils/api";
import authenticateUser from "../../../utils/authenticateUser";

// Get all inboxes
export const get: APIRoute = async ({ params, request }) => {
  try {
    const { user: currentUser, isAuthorized } = await authenticateUser(request);
    if (!isAuthorized) {
      return errorResponse(null, 401, "Unauthorized request");
    }

    const queryParams = getQueryParams(request.url);
    const limit: string | null = queryParams.get("limit") || "10";
    const cursor: string | null = queryParams.get("cursor");

    const db: Db = await getDb();
    const rooms = db.collection("rooms");
    const allRooms = await rooms
      .find(
        {},
        {
          sort: {
            _id: -1,
          },
          limit: Number(limit),
        }
      )
      .toArray();

    return successResponse(allRooms, {
      cursor,
      size: limit,
    });
  } catch (err) {
    console.error(err);
    return errorResponse(err);
  }
};
