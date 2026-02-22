import { Injectable, Logger } from '@nestjs/common';
import { ISessionStore } from './session-store.interface';
import { SessionDataDto } from '../dto/session.dto';

interface InternalSessionRecord {
   data: SessionDataDto;
   expireAt: Date;
}

@Injectable()
export class LocalSessionStore extends ISessionStore {
   private readonly logger = new Logger(LocalSessionStore.name);
   private readonly store = new Map<string, InternalSessionRecord>();

   set(sid: string, session: SessionDataDto): Promise<void> {
      const expiresAt = this.calculateExpiry(session);
      this.store.set(sid, { data: session, expireAt: expiresAt });
      return Promise.resolve();
   }

   get(sid: string): Promise<SessionDataDto | null> {
      const rec = this.store.get(sid);
      if (!rec) return Promise.resolve(null);
      if (rec.expireAt < new Date()) {
         this.store.delete(sid);
         return Promise.resolve(null);
      }
      return Promise.resolve(rec.data);
   }

   destroy(sid: string): Promise<void> {
      this.store.delete(sid);
      return Promise.resolve();
   }

   touch(sid: string, session: SessionDataDto): Promise<void> {
      const rec = this.store.get(sid);
      if (!rec) return Promise.resolve();
      const expiresAt = this.calculateExpiry(session);
      rec.expireAt = expiresAt;
      this.store.set(sid, rec);
      return Promise.resolve();
   }

   clearExpired(): Promise<void> {
      const now = new Date();
      let removed = 0;
      for (const [sid, rec] of this.store.entries()) {
         if (rec.expireAt < now) {
            this.store.delete(sid);
            removed++;
         }
      }
      if (removed > 0) this.logger.log(`Cleared ${removed} expired sessions`);
      return Promise.resolve();
   }

   private calculateExpiry(session: SessionDataDto): Date {
      const cookieExpires = session.cookie?.expires;

      if (cookieExpires instanceof Date) {
         return cookieExpires;
      }

      if (typeof cookieExpires === 'string') {
         return new Date(cookieExpires);
      }

      return new Date(Date.now() + 24 * 60 * 60 * 1000);
   }
}
