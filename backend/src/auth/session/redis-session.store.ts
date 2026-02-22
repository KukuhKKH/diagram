import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ISessionStore } from './session-store.interface';
import { SessionDataDto } from '../dto/session.dto';

@Injectable()
export class RedisSessionStore extends ISessionStore {
   private readonly logger = new Logger(RedisSessionStore.name);

   private redis: {
      set: (
         key: string,
         value: string,
         mode?: string,
         ttl?: number,
      ) => Promise<unknown>;
      get: (key: string) => Promise<string | null>;
      del: (key: string) => Promise<number>;
      expire: (key: string, ttl: number) => Promise<number>;
   } | null = null;

   constructor(private readonly configService: ConfigService) {
      super();

      const redisUrl = this.configService.get<string>(
         'REDIS_URL',
         'redis://localhost:6379',
      );

      void import('ioredis')
         .then((mod) => {
            const maybeCtor =
               (mod as unknown as { default?: unknown }).default ??
               (mod as unknown);

            const Ctor = maybeCtor as { new (url: string): unknown };

            const clientObj = new Ctor(redisUrl) as {
               set: (
                  key: string,
                  value: string,
                  mode?: string,
                  ttl?: number,
               ) => Promise<unknown>;
               get: (key: string) => Promise<string | null>;
               del: (key: string) => Promise<number>;
               expire: (key: string, ttl: number) => Promise<number>;
            };

            this.redis = clientObj;
            this.logger.log('Connected to Redis for session storage');
         })
         .catch((err: unknown) => {
            this.redis = null;
            this.logger.warn(
               'ioredis not installed or failed to initialize. RedisSessionStore disabled.',
            );
            this.logger.debug(String(err));
         });
   }

   async set(sid: string, session: SessionDataDto): Promise<void> {
      if (!this.redis) {
         throw new Error(
            'Redis client not initialized. Install/configure ioredis',
         );
      }

      const ttl = this.calculateTTL(session);
      const key = `session:${sid}`;
      const value = JSON.stringify(session.toJSON());

      await this.redis.set(key, value, 'EX', ttl);
   }

   async get(sid: string): Promise<SessionDataDto | null> {
      if (!this.redis) {
         throw new Error(
            'Redis client not initialized. Install/configure ioredis',
         );
      }

      const key = `session:${sid}`;
      const data = await this.redis.get(key);
      if (!data) return null;

      try {
         const parsed = JSON.parse(data) as Record<string, unknown>;
         return SessionDataDto.fromJSON(parsed);
      } catch (err) {
         this.logger.error('Failed to parse session JSON from Redis', err);
         return null;
      }
   }

   async destroy(sid: string): Promise<void> {
      if (!this.redis) {
         throw new Error(
            'Redis client not initialized. Install/configure ioredis',
         );
      }

      const key = `session:${sid}`;
      await this.redis.del(key);
   }

   async touch(sid: string, session: SessionDataDto): Promise<void> {
      if (!this.redis) {
         throw new Error(
            'Redis client not initialized. Install/configure ioredis',
         );
      }

      const ttl = this.calculateTTL(session);
      const key = `session:${sid}`;
      await this.redis.expire(key, ttl);
   }

   async clearExpired(): Promise<void> {
      this.logger.log('Redis automatically handles expired sessions');
      return Promise.resolve();
   }

   private calculateTTL(session: SessionDataDto): number {
      const cookieExpires = session.cookie?.expires;
      const expiresTime = cookieExpires
         ? new Date(cookieExpires).getTime()
         : Date.now() + 24 * 60 * 60 * 1000;

      const ttlMs = expiresTime - Date.now();
      return Math.max(Math.floor(ttlMs / 1000), 60); // At least 60 seconds
   }
}
