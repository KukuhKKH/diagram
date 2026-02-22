import {
   Injectable,
   NestMiddleware,
   Logger,
   ForbiddenException,
} from '@nestjs/common';
import type { Request, Response, NextFunction } from 'express';

@Injectable()
export class CsrfProtectionMiddleware implements NestMiddleware {
   private readonly logger = new Logger(CsrfProtectionMiddleware.name);

   use(req: Request, res: Response, next: NextFunction) {
      if (['GET', 'HEAD', 'OPTIONS'].includes(req.method)) {
         return next();
      }

      if (req.path === '/auth/callback') {
         return next();
      }

      try {
         const token =
            (req.headers['x-csrf-token'] as string | undefined) ||
            ((req.body as Record<string, unknown>)?.['_csrf'] as
               | string
               | undefined);

         if (!token) {
            this.logger.warn(
               `Missing CSRF token for ${req.method} ${req.path}`,
            );
            throw new ForbiddenException('CSRF token required');
         }

         // TODO: Implement actual CSRF validation logic
         // For now, just pass through

         next();
      } catch (error) {
         this.logger.error(`CSRF check failed: ${error}`);
         throw new ForbiddenException('CSRF validation failed');
      }
   }
}

export function generateCsrfToken(): string {
   // TODO: Implement token generation from double-CSRF library
   return 'stub-token';
}
