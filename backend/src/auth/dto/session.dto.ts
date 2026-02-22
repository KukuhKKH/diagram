export class SessionDataDto {
   cookie: SessionCookieDto;
   user?: SessionUserDto;
   logtoAccessToken?: string;
   logtoRefreshToken?: string;
   issuedAt?: number;
   expiresAt?: number;

   constructor(data: Partial<SessionDataDto>) {
      this.cookie = data.cookie ?? new SessionCookieDto();
      this.user = data.user;
      this.logtoAccessToken = data.logtoAccessToken;
      this.logtoRefreshToken = data.logtoRefreshToken;
      this.issuedAt = data.issuedAt;
      this.expiresAt = data.expiresAt;
   }

   toJSON(): Record<string, unknown> {
      return {
         cookie: this.cookie,
         user: this.user,
         logtoAccessToken: this.logtoAccessToken,
         logtoRefreshToken: this.logtoRefreshToken,
         issuedAt: this.issuedAt,
         expiresAt: this.expiresAt,
      };
   }

   static fromJSON(data: Record<string, unknown>): SessionDataDto {
      return new SessionDataDto({
         cookie: data.cookie as SessionCookieDto,
         user: data.user as SessionUserDto | undefined,
         logtoAccessToken: data.logtoAccessToken as string | undefined,
         logtoRefreshToken: data.logtoRefreshToken as string | undefined,
         issuedAt: data.issuedAt as number | undefined,
         expiresAt: data.expiresAt as number | undefined,
      });
   }
}

export class SessionCookieDto {
   expires?: Date;
   originalMaxAge?: number;
   httpOnly?: boolean;
   secure?: boolean;
   sameSite?: 'strict' | 'lax' | 'none' | boolean;
   path?: string;

   constructor(data?: Partial<SessionCookieDto>) {
      this.expires = data?.expires;
      this.originalMaxAge = data?.originalMaxAge;
      this.httpOnly = data?.httpOnly;
      this.secure = data?.secure;
      this.sameSite = data?.sameSite;
      this.path = data?.path;
   }
}

export class SessionUserDto {
   id: string;
   logtoUserId: string;
   email?: string;
   name?: string;
   avatarUrl?: string;

   constructor(data: SessionUserDto) {
      this.id = data.id;
      this.logtoUserId = data.logtoUserId;
      this.email = data.email;
      this.name = data.name;
      this.avatarUrl = data.avatarUrl;
   }
}

export interface SessionStoreResult {
   sid: string;
   session: SessionDataDto | null;
}
