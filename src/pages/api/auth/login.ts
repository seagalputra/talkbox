import type { APIRoute } from "astro";

export const post: APIRoute = ({ params, request }) => {
  return {
    body: JSON.stringify({}),
  };
};
