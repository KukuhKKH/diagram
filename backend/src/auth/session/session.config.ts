import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class SessionConfigService {
   private readonly logger = new Logger(SessionConfigService.name);

   constructor(private readonly configService: ConfigService) {}

   /**
    * Get session store type from environment.
    * @returns 'local' | 'redis'
    */
   getStoreType(): 'local' | 'redis' {
      const storeType = this.configService.get<string>(
         'SESSION_STORE_TYPE',
         'local',
      );

      if (storeType !== 'local' && storeType !== 'redis') {
         this.logger.warn(
            `Invalid SESSION_STORE_TYPE: ${storeType}. Defaulting to 'local'.`,
         );
         return 'local';
      }

      return storeType;
   }

   /**
    * Get session secret for signing cookies.
    */
   getSecret(): string {
      return this.configService.getOrThrow<string>('SESSION_SECRET');
   }

   /**
    * Get session cookie name.
    */
   getCookieName(): string {
      return this.configService.get<string>(
         'SESSION_COOKIE_NAME',
         'diagram_session',
      );
   }

   /**
    * Get session max age in milliseconds.
    */
   getMaxAge(): number {
      const hours = this.configService.get<number>('SESSION_MAX_AGE_HOURS', 24);
      return hours * 60 * 60 * 1000;
   }

   /**
    * Check if running in production.
    */
   isProduction(): boolean {
      return this.configService.get<string>('NODE_ENV') === 'production';
   }

   /**
    * Get Redis URL if using Redis store.
    */
   getRedisUrl(): string {
      return this.configService.get<string>(
         'REDIS_URL',
         'redis://localhost:6379',
      );
   }
}
