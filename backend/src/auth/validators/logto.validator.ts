import type { LogtoUserProfile } from '../types/session.types';

export async function validateLogtoProfile(
   raw: unknown,
): Promise<LogtoUserProfile> {
   try {
      /* eslint-disable @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call */
      const mod = await import('zod');
      const z = (mod as unknown as { z?: unknown }).z ?? (mod as unknown);

      const schema = (z as any).object({
         sub: (z as any).string().min(1),
         email: (z as any).string().email().optional(),
         name: (z as any).string().optional(),
         picture: (z as any).string().optional(),
         identities: (z as any).record((z as any).unknown()).optional(),
      });

      const parsed = schema.parse(raw);
      const result: LogtoUserProfile = {
         id: parsed.sub,
         email: parsed.email,
         name: parsed.name,
         picture: parsed.picture,
         identities: parsed.identities,
      };
      /* eslint-enable @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call */

      return result;
   } catch {
      // Fallback: manual validation
      if (typeof raw !== 'object' || raw === null) {
         throw new Error('Invalid user profile: not an object');
      }

      const obj = raw as Record<string, unknown>;
      const sub = obj.sub;
      if (typeof sub !== 'string' || sub.length === 0) {
         throw new Error('Invalid user profile from Logto: missing sub');
      }

      const email = typeof obj.email === 'string' ? obj.email : undefined;
      const name = typeof obj.name === 'string' ? obj.name : undefined;
      const picture = typeof obj.picture === 'string' ? obj.picture : undefined;
      const identities =
         obj.identities && typeof obj.identities === 'object'
            ? (obj.identities as Record<string, unknown>)
            : undefined;

      return {
         id: sub,
         email,
         name,
         picture,
         identities,
      };
   }
}
