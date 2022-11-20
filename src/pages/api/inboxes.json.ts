export async function get({ params, request }: { params: any; request: any }) {
  return {
    body: JSON.stringify({
      name: "Astro",
      url: "https://astro.build",
    }),
  };
}
