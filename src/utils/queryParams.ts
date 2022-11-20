function getQueryParams(url: string) {
  const urlObject = new URL(url);
  return new URLSearchParams(urlObject.search);
}

export { getQueryParams };
