/**
 * Session types for Logto authentication via BFF pattern.
 * Represents authenticated user data stored in httpOnly session cookie.
 */

export interface SessionUser {
   id: string;
   logtoUserId: string;
   email?: string;
   name?: string;
   avatarUrl?: string;
}

export interface AuthSession {
   user: SessionUser | null;
   logtoAccessToken?: string;
   logtoRefreshToken?: string;
   issuedAt: number;
   expiresAt: number;
}

export interface LogtoUserProfile {
   id: string;
   email?: string;
   name?: string;
   picture?: string;
   identities?: Record<string, unknown>;
}

export interface OAuthTokens {
   accessToken: string;
   refreshToken?: string;
   idToken?: string;
   expiresIn?: number;
}
