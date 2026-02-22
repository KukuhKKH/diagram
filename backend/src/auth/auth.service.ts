import { Injectable, Logger, UnauthorizedException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import type { SessionUser, LogtoUserProfile } from './types/session.types';

/**
 * Auth Service: Orchestrates authentication logic.
 * Handles user persistence, session management, and Logto integration.
 *
 * SOLID Principle - Dependency Inversion:
 * Depends on PrismaService abstraction, not direct database access.
 * Handles business logic separate from HTTP concerns.
 */
@Injectable()
export class AuthService {
   private readonly logger = new Logger(AuthService.name);

   constructor(private readonly prisma: PrismaService) {}

   /**
    * Find or create user from Logto profile.
    * Synchronizes Logto user data with local database.
    *
    * Flow:
    * 1. Try to find user by logtoUserId
    * 2. If not found, create new user
    * 3. Update profile fields if changed
    * 4. Return as SessionUser
    */
   async findOrCreateUser(
      logtoProfile: LogtoUserProfile,
   ): Promise<SessionUser> {
      try {
         let user = await this.prisma.user.findUnique({
            where: { logtoUserId: logtoProfile.id },
         });

         if (!user) {
            // Create new user
            user = await this.prisma.user.create({
               data: {
                  logtoUserId: logtoProfile.id,
                  email: logtoProfile.email,
                  name: logtoProfile.name,
                  avatarUrl: logtoProfile.picture,
               },
            });

            this.logger.log(`Created new user: ${user.id}`);
         } else {
            // Update if profile changed
            const updateData: {
               email?: string;
               name?: string;
               avatarUrl?: string;
            } = {};

            if (logtoProfile.email && user.email !== logtoProfile.email) {
               updateData.email = logtoProfile.email;
            }
            if (logtoProfile.name && user.name !== logtoProfile.name) {
               updateData.name = logtoProfile.name;
            }
            if (
               logtoProfile.picture &&
               user.avatarUrl !== logtoProfile.picture
            ) {
               updateData.avatarUrl = logtoProfile.picture;
            }

            if (Object.keys(updateData).length > 0) {
               user = await this.prisma.user.update({
                  where: { id: user.id },
                  data: updateData,
               });

               this.logger.log(`Updated user profile: ${user.id}`);
            }
         }

         return this.toSessionUser(user);
      } catch (error) {
         this.logger.error('Failed to find or create user', error);
         throw new UnauthorizedException('Failed to process user data');
      }
   }

   /**
    * Convert database User to SessionUser.
    * Removes sensitive fields, keeps only what's needed for session.
    */
   private toSessionUser(user: {
      id: string;
      logtoUserId: string;
      email: string | null;
      name: string | null;
      avatarUrl: string | null;
   }): SessionUser {
      return {
         id: user.id,
         logtoUserId: user.logtoUserId,
         email: user.email ?? undefined,
         name: user.name ?? undefined,
         avatarUrl: user.avatarUrl ?? undefined,
      };
   }

   /**
    * Get user by ID.
    * Useful for session validation and profile endpoints.
    */
   async getUserById(userId: string): Promise<SessionUser | null> {
      try {
         const user = await this.prisma.user.findUnique({
            where: { id: userId },
         });

         return user ? this.toSessionUser(user) : null;
      } catch (error) {
         this.logger.error(`Failed to get user ${userId}`, error);
         return null;
      }
   }

   /**
    * Logout user: clear session tokens.
    * In this BFF pattern, actual logout is handled by clearing cookies.
    * This method can be extended for additional cleanup (token revocation, etc).
    */
   logout(userId: string): void {
      try {
         // Optional: Update last logout timestamp, revoke tokens, etc.
         this.logger.log(`User logged out: ${userId}`);
         // In future, can revoke refresh tokens here
      } catch (error) {
         this.logger.error(`Logout failed for user ${userId}`, error);
      }
   }
}
