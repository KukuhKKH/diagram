import type { SessionUser } from './session.types';

/**
 * Passport OAuth profile and tokens type.
 * Used in Logto strategy callback.
 */
export interface PassportOAuthUser {
   userProfile: {
      id: string;
      email?: string;
      name?: string;
      picture?: string;
      identities?: Record<string, unknown>;
   };
   tokens: {
      accessToken: string;
      refreshToken?: string;
      idToken?: string;
      expiresIn?: number;
   };
}

/**
 * Extend Express Request type to include session and user properties.
 * Used for proper TypeScript typing in middleware, decorators, and guards.
 */
declare global {
   namespace Express {
      interface Request {
         session: {
            user?: SessionUser;
            logtoAccessToken?: string;
            logtoRefreshToken?: string;
            issuedAt?: number;
            expiresAt?: number;
            save?: (callback?: (err?: Error) => void) => void;
            destroy?: (callback?: (err?: Error) => void) => void;
         } & Record<string, unknown>;
         user?: SessionUser | PassportOAuthUser;
      }
   }
}

export {};
