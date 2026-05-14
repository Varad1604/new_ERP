import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value
  const { pathname } = request.nextUrl

  // Allow access to login, API routes, and static files
  if (pathname.startsWith('/login') || pathname.startsWith('/_next') || pathname.startsWith('/api')) {
    return NextResponse.next()
  }

  // Redirect to login if no token is found (server-side check)
  // For now, since Zustand manages the token locally, we let the client handle redirects for this phase
  // In a full production app, you would verify an HttpOnly JWT cookie here.

  if (pathname === '/') {
     return NextResponse.redirect(new URL('/login', request.url))
  }

  return NextResponse.next()
}
