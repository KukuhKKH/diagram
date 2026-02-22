import { SessionDataDto } from '../dto/session.dto';

export abstract class ISessionStore {
   /**
    * Store session data.
    * @param sid - Session ID
    * @param session - Session data DTO
    * @returns Promise that resolves when session is stored
    */
   abstract set(sid: string, session: SessionDataDto): Promise<void>;

   /**
    * Retrieve session data by ID.
    * @param sid - Session ID
    * @returns Promise with session data or null if not found
    */
   abstract get(sid: string): Promise<SessionDataDto | null>;

   /**
    * Delete session by ID.
    * @param sid - Session ID
    * @returns Promise that resolves when session is deleted
    */
   abstract destroy(sid: string): Promise<void>;

   /**
    * Touch session to update expiry without modifying data.
    * @param sid - Session ID
    * @param session - Session data DTO
    * @returns Promise that resolves when session is touched
    */
   abstract touch(sid: string, session: SessionDataDto): Promise<void>;

   /**
    * Clear all expired sessions.
    * Should be called periodically by a cron job.
    * @returns Promise that resolves when expired sessions are cleared
    */
   abstract clearExpired(): Promise<void>;
}
