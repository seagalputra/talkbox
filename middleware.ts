import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  const hasTalkboxCookie = request.cookies.has("talkbox");
  if (pathname === "/") {
    if (hasTalkboxCookie) {
      return NextResponse.redirect(new URL("/inboxes", request.url));
    }
  } else {
    if (!hasTalkboxCookie) {
      return NextResponse.redirect(new URL("/", request.url));
    } else {
      return NextResponse.next();
    }
  }
}

export const config = {
  matcher: ["/inboxes/:path*", "/"],
};
