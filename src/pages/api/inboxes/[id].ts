import type { APIRoute } from "astro";

// Get inboxes specific by id
export const get: APIRoute = ({ params, request }) => {
  const { id } = params;

  return {
    body: JSON.stringify({
      name: "Inboxes endpoint",
      id,
    }),
  };
};
