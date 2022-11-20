import type { Db } from "mongodb";
import type { APIRoute } from "astro";
import { getDb } from "../../../assets/mongo";
import { getQueryParams } from "../../../utils/queryParams";

// Get all inboxes
export const get: APIRoute = async ({ params, request }) => {
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

  return new Response(
    JSON.stringify({
      status: "success",
      meta: {
        cursor,
        size: limit,
      },
      data: allRooms,
    }),
    {
      status: 200,
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
};
