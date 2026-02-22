import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import type { Request } from 'express';
import type { SessionUser } from '../types/session.types';
import '../types/express.d';

export const CurrentUser = createParamDecorator(
   (data: unknown, ctx: ExecutionContext): SessionUser | null => {
      const request = ctx.switchToHttp().getRequest<Request>();
      const session = request.session;

      if (!session || !session.user) {
         return null;
      }

      return session.user;
   },
);

export const RequireAuth = createParamDecorator(
   (data: unknown, ctx: ExecutionContext): SessionUser => {
      const request = ctx.switchToHttp().getRequest<Request>();
      const session = request.session;

      if (!session || !session.user) {
         throw new Error('Unauthorized: no session found');
      }

      return session.user;
   },
);
