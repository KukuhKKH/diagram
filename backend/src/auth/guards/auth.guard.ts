import {
   Injectable,
   CanActivate,
   ExecutionContext,
   UnauthorizedException,
} from '@nestjs/common';
import type { Request } from 'express';
import type { SessionUser } from '../types/session.types';
import '../types/express.d';

@Injectable()
export class AuthGuard implements CanActivate {
   canActivate(context: ExecutionContext): boolean {
      const request = context.switchToHttp().getRequest<Request>();
      const sessionUser = this.getSessionUser(request);

      if (!sessionUser || !sessionUser.id) {
         throw new UnauthorizedException('Not authenticated');
      }

      request.user = sessionUser;

      return true;
   }

   private getSessionUser(request: Request): SessionUser | null {
      const session = request.session;

      if (!session || !session.user) {
         return null;
      }

      return session.user;
   }
}

@Injectable()
export class RoleGuard implements CanActivate {
   constructor(private allowedRoles: string[]) {}

   canActivate(context: ExecutionContext): boolean {
      const request = context.switchToHttp().getRequest<Request>();
      const sessionUser = request.session?.user ?? null;

      if (!sessionUser) {
         throw new UnauthorizedException('Not authenticated');
      }

      // TODO: Implement role checking when roles are added to user schema

      return true;
   }
}
