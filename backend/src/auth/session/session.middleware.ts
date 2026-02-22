import { Injectable, NestMiddleware, Logger } from '@nestjs/common';
import type { Request, Response, NextFunction } from 'express';
import * as session from 'express-session';
import type { Store } from 'express-session';
import { ISessionStore } from './session-store.interface';
import { SessionConfigService } from './session.config';
import { SessionDataDto } from '../dto/session.dto';

/* eslint-disable @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call */
@Injectable()
export class SessionMiddleware implements NestMiddleware {
   private readonly logger = new Logger(SessionMiddleware.name);

   constructor(
      private readonly sessionStore: ISessionStore,
      private readonly sessionConfig: SessionConfigService,
   ) {}

   use(req: Request, res: Response, next: NextFunction): void {
      const sessionMiddleware = session({
         store: this.createExpressSessionStore(),
         secret: this.sessionConfig.getSecret(),
         resave: false,
         saveUninitialized: false,
         cookie: {
            httpOnly: true,
            secure: this.sessionConfig.isProduction(),
            sameSite: 'lax',
            maxAge: this.sessionConfig.getMaxAge(),
            path: '/',
         },
         name: this.sessionConfig.getCookieName(),
      });

      (
         sessionMiddleware as unknown as (
            req: Request,
            res: Response,
            next: NextFunction,
         ) => void
      )(req, res, next);
   }

   private createExpressSessionStore(): Partial<Store> {
      const store = this.sessionStore;

      return {
         get: (
            sid: string,
            callback: (
               err?: Error | null,
               session?: Record<string, unknown> | null,
            ) => void,
         ): void => {
            void store
               .get(sid)
               .then((sessionDto: SessionDataDto | null) => {
                  if (!sessionDto) {
                     callback(null, null);
                     return;
                  }
                  const sessionData = sessionDto.toJSON();
                  callback(null, sessionData);
               })
               .catch((err: unknown) => {
                  if (err instanceof Error) callback(err, null);
                  else callback(new Error(String(err)), null);
               });
         },
         set: (
            sid: string,
            session: Record<string, unknown>,
            callback?: (err?: Error | null) => void,
         ): void => {
            const sessionDto = SessionDataDto.fromJSON(session);

            void store
               .set(sid, sessionDto)
               .then(() => {
                  if (callback) callback();
               })
               .catch((err: unknown) => {
                  if (callback) {
                     if (err instanceof Error) callback(err);
                     else callback(new Error(String(err)));
                  }
               });
         },
         destroy: (
            sid: string,
            callback?: (err?: Error | null) => void,
         ): void => {
            void store
               .destroy(sid)
               .then(() => {
                  if (callback) callback();
               })
               .catch((err: unknown) => {
                  if (callback) {
                     if (err instanceof Error) callback(err);
                     else callback(new Error(String(err)));
                  }
               });
         },
         touch: (
            sid: string,
            session: Record<string, unknown>,
            callback?: (err?: Error | null) => void,
         ): void => {
            const sessionDto = SessionDataDto.fromJSON(session);

            void store
               .touch(sid, sessionDto)
               .then(() => {
                  if (callback) callback();
               })
               .catch((err: unknown) => {
                  if (callback) {
                     if (err instanceof Error) callback(err);
                     else callback(new Error(String(err)));
                  }
               });
         },
      };
   }
}
