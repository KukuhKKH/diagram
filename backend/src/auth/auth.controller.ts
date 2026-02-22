import {
   Controller,
   Get,
   Post,
   Res,
   UseGuards,
   UnauthorizedException,
   Logger,
} from '@nestjs/common';
import type { Response, Request } from 'express';
import { AuthGuard as PassportAuthGuard } from '@nestjs/passport';
import { AuthGuard } from './guards/auth.guard';
import { CurrentUser } from './decorators/current-user.decorator';
import { AuthService } from './auth.service';
import type { SessionUser } from './types/session.types';
import type { PassportOAuthUser } from './types/express.d';
import { AuthStatusDto, UserResponseDto } from './dto/user.dto';
import './types/express.d';

@Controller('auth')
export class AuthController {
   private readonly logger = new Logger(AuthController.name);

   constructor(private readonly authService: AuthService) {}

   @Get('login')
   @UseGuards(PassportAuthGuard('logto'))
   login(): void {}

   @Get('callback')
   @UseGuards(PassportAuthGuard('logto'))
   callback(req: Request, @Res() res: Response) {
      const oauthUser = req.user as PassportOAuthUser | undefined;
      if (!oauthUser) {
         throw new UnauthorizedException('OAuth callback failed');
      }

      const { userProfile, tokens } = oauthUser;

      req.session.user = {
         id: userProfile.id,
         logtoUserId: userProfile.id,
         email: userProfile.email,
         name: userProfile.name,
         avatarUrl: userProfile.picture,
      };

      req.session.logtoAccessToken = tokens.accessToken;
      req.session.logtoRefreshToken = tokens.refreshToken;
      req.session.issuedAt = Date.now();
      req.session.expiresAt = Date.now() + (tokens.expiresIn || 3600) * 1000;

      // Save session
      const sessionSave = req.session?.save as
         | ((cb: (err?: Error) => void) => void)
         | undefined;
      if (sessionSave) {
         sessionSave((err?: Error) => {
            if (err) {
               this.logger.error('Failed to save session', err);
               res.redirect('/auth/error');
            } else {
               // Redirect to frontend (adjust URL as needed)
               res.redirect(
                  process.env.FRONTEND_URL || 'http://localhost:5173',
               );
            }
         });
      } else {
         // Fallback if save is not available
         res.redirect(process.env.FRONTEND_URL || 'http://localhost:5173');
      }
   }

   @Get('status')
   getStatus(@CurrentUser() user: SessionUser | null): AuthStatusDto {
      return {
         isAuthenticated: !!user,
         user: user
            ? {
                 id: user.id,
                 logtoUserId: user.logtoUserId,
                 email: user.email,
                 name: user.name,
                 avatarUrl: user.avatarUrl,
                 createdAt: new Date(),
                 updatedAt: new Date(),
              }
            : undefined,
      };
   }

   @Get('profile')
   @UseGuards(AuthGuard)
   async getProfile(
      @CurrentUser() user: SessionUser,
   ): Promise<UserResponseDto> {
      // CurrentUser already guaranteed to exist via AuthGuard
      const dbUser = await this.authService.getUserById(user.id);

      if (!dbUser) {
         throw new UnauthorizedException('User not found');
      }

      return {
         id: dbUser.id,
         logtoUserId: dbUser.logtoUserId,
         email: dbUser.email,
         name: dbUser.name,
         avatarUrl: dbUser.avatarUrl,
         createdAt: new Date(),
         updatedAt: new Date(),
      };
   }

   @Post('logout')
   @UseGuards(AuthGuard)
   logout(req: Request, @Res() res: Response): void {
      const userId = req.session?.user?.id;

      const sessionDestroy = req.session?.destroy as
         | ((cb: (err?: Error) => void) => void)
         | undefined;
      if (sessionDestroy) {
         sessionDestroy((err?: Error) => {
            if (err) {
               this.logger.error('Session destroy failed', err);
               res.status(500).json({ message: 'Logout failed' });
               return;
            }

            if (userId) {
               try {
                  this.authService.logout(userId);
               } catch (e) {
                  this.logger.error('Post-logout cleanup failed', e as Error);
               }
            }

            res.clearCookie('diagram_session');
            res.json({ message: 'Logged out successfully' });
         });
      } else {
         this.logger.error('Session destroy method not available');
         res.status(500).json({ message: 'Logout failed' });
      }
   }

   /**
    * TODO: Implement token refresh logic
    */
   @Post('refresh')
   @UseGuards(AuthGuard)
   refresh(req: Request, @Res() res: Response) {
      // TODO: Implement refresh token logic
      res.json({ message: 'Token refresh not yet implemented' });
   }

   @Get('error')
   error(@Res() res: Response) {
      res.status(401).json({ message: 'Authentication failed' });
   }
}
